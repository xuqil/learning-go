//go:build v7

package web

import (
	"fmt"
	"regexp"
	"strings"
)

// router 路由
type router struct {
	// trees 路由树（森林）
	// key 为 HTTP 方法
	// value 为 路由树
	trees map[string]*node
}

func newRouter() router {
	return router{trees: map[string]*node{}}
}

// addRoute 注册路由。
// method 是 HTTP 方法
// - 已经注册了的路由，无法被覆盖。例如 /user/home 注册两次，会冲突
// - path 必须以 / 开始并且结尾不能有 /，中间也不允许有连续的 /
// - 不能在同一个位置注册不同的参数路由，例如 /user/:id 和 /user/:name 冲突
// - 不能在同一个位置同时注册通配符路由和参数路由，例如 /user/:id 和 /user/* 冲突
// - 同名路径参数，在路由匹配的时候，值会被覆盖。例如 /user/:id/abc/:id，那么 /user/123/abc/456 最终 id = 456
func (r *router) addRoute(method string, path string, handler HandleFunc) {
	if path == "" {
		panic("route: 路由是空字符串")
	}

	if path[0] != '/' {
		panic("route: 路由必须以 / 开头")
	}

	if path != "/" && path[len(path)-1] == '/' {
		panic("route: 路由不能以 / 结尾")
	}

	root, ok := r.trees[method]
	if !ok {
		root = &node{path: "/"}
		root.typ = nodeTypeStatic
		r.trees[method] = root
	}

	// 特殊路由处理
	if path == "/" {
		if root.handler != nil {
			panic("route: 路由冲突[/]")
		}
		root.typ = nodeTypeStatic
		root.handler = handler
		return
	}

	// 切割路由
	for _, seg := range strings.Split(path[1:], "/") {
		if seg == "" {
			panic(fmt.Sprintf("route: 非法路由，路由不能有连续的 /, [%s]", path))
		}
		root = root.childOrCreate(seg)
	}

	if root.handler != nil {
		panic(fmt.Sprintf("route: 路由冲突[%s]", path))
	}
	root.handler = handler
}

// findRoute 路由查找
// 找到了 node，但 node 的 handler 不一定有
func (r *router) findRoute(method string, path string) (*matchInfo, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	// 特殊路径
	if path == "/" {
		return &matchInfo{n: root}, true
	}
	// 切割 path 并匹配路由
	mi := &matchInfo{}
	for _, seg := range strings.Split(strings.Trim(path, "/"), "/") {
		var matchParam bool
		root, matchParam, ok = root.childOf(seg)
		if !ok {
			return nil, ok
		}
		if matchParam {
			mi.addValue(root.path[1:], seg)
		}
	}
	mi.n = root
	return mi, true
}

type nodeType int

const (
	//	静态路由
	nodeTypeStatic = iota
	// 正则路由
	nodeTypeReg
	// 路径参数路由
	nodeTypeParam
	// 通配符路由
	nodeTypeAny
)

// node 代表路由树的节点
// 路由树的匹配顺序是：
// 1. 静态完全匹配
// 2. 正则匹配，形式 :param_name(reg_expr)
// 3. 路径参数匹配：形式 :param_name
// 4. 通配符匹配：*
// 这是不回溯匹配
type node struct {
	// 节点类型，默认是 nodeTypeStatic
	typ nodeType

	// path 路由
	path string

	// children 子节点
	// key 为 path
	// value 是 子节点
	children map[string]*node

	// starChild 通配符子节点
	starChild *node

	// paramChild 参数路径子节点
	paramChild *node
	// 正则路由和参数路由都会使用这个字段
	paramName string

	// 正则表达式
	regChild *node
	regExpr  *regexp.Regexp

	// handler 业务逻辑
	handler HandleFunc
}

// childOrCreate 查找子节点，
// 首先会判断 path 是不是通配符路径
// 其次判断 path 是不是参数路径，即以 : 开头的路径
// 最后会从 children 里面查找，
// 如果没有找到，那么会创建一个新的节点，并且保存在 node 里面
func (n *node) childOrCreate(path string) *node {
	// 通配符路由节点
	if path == "*" {
		return n.childOrCreateStar(path)
	}

	// 参数路径节点
	if path[0] == ':' {
		return n.childOrCreateParamReg(path)
	}

	// 如果 children 没有初始化，则进行初始化
	if n.children == nil {
		n.children = map[string]*node{}
	}
	child, ok := n.children[path]
	if !ok {
		// 没有的话需要创建一个 node
		child = &node{path: path}
		child.typ = nodeTypeStatic
		n.children[path] = child
	}
	return child
}

// childOrCreateStar 创建通配符路由
func (n *node) childOrCreateStar(path string) *node {
	if n.regChild != nil {
		panic(fmt.Sprintf("route: 非法路由，已有正则路由。不允许同时注册通配符路由和正则路由 [%s]", path))
	}
	if n.paramChild != nil {
		panic(fmt.Sprintf("route: 非法路由，已有路径参数路由。不允许同时注册通配符路由和参数路由 [%s]", path))
	}
	if n.starChild == nil {
		starChild := &node{path: path}
		starChild.typ = nodeTypeAny
		n.starChild = starChild
	}
	return n.starChild
}

// childOrCreateParamReg 创建参数路径路由或者正则路由
func (n *node) childOrCreateParamReg(path string) *node {
	param, regStr, ok := parseParam(path)
	if !ok {
		if n.regChild != nil {
			panic(fmt.Sprintf("route: 非法路由，已有正则路由。不允许同时注册正则路由和参数路由 [%s]", path))
		}
		if n.starChild != nil {
			panic(fmt.Sprintf("route: 非法路由，已有通配符路由。不允许同时注册通配符路由和参数路由 [%s]", path))
		}
		if n.paramChild != nil {
			if n.paramChild.path != path {
				panic(fmt.Sprintf("route: 路由冲突，参数路由冲突，已有 %s，新注册 %s", n.paramChild.path, path))
			}
		} else {
			n.paramChild = &node{path: path}
			n.paramChild.typ = nodeTypeParam
		}
		return n.paramChild
	}
	if n.starChild != nil {
		panic(fmt.Sprintf("route: 非法路由，已有通配符路由。不允许同时注册通配符路由和正则路由 [%s]", path))
	}
	if n.paramChild != nil {
		panic(fmt.Sprintf("route: 非法路由，已有路径参数路由。不允许同时注册正则路由和参数路由 [%s]", path))
	}
	if n.regChild != nil {
		if n.regExpr.String() != path {
			panic(fmt.Sprintf("route: 非法路由，已有正则路由。不允许同时注册正则路由和参数路由 [%s]", path))
		}
	} else {
		regChild := &node{path: path}
		regChild.typ = nodeTypeReg
		regChild.paramName = param
		regChild.regExpr = regexp.MustCompile(regStr)
		n.regChild = regChild
	}
	return n.regChild
}

// childOf 返回子节点
// 第一个返回值 *node 是命中的节点
// 第二个返回值 bool 代表是否命中
func (n *node) childOf(path string) (*node, bool, bool) {
	if n.children == nil {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		if n.starChild != nil {
			return n.starChild, false, true
		}
		return n, false, n.path == "*"
	}
	child, ok := n.children[path]
	if !ok {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		if n.starChild != nil {
			return n.starChild, false, true
		}
		return n, false, n.path == "*"
	}
	return child, false, ok
}

// matchInfo 匹配的子节点和路径参数
type matchInfo struct {
	n *node

	pathParams map[string]string
}

func (m *matchInfo) addValue(key, value string) {
	if m.pathParams == nil {
		m.pathParams = map[string]string{}
	}
	m.pathParams[key] = value
}

// parseParam 解析路径参数
// 第一个参数返回的是路径参数
// 第二个返回的是正则表达式
// 第三个代表是否有正在表达式
func parseParam(path string) (string, string, bool) {
	path = path[1:]
	res := strings.SplitN(path, "(", 2)
	if len(res) == 2 {
		param := res[0]
		if strings.HasSuffix(path, ")") {
			return param, path[len(param):], true
		}
	}
	return path, "", false
}
