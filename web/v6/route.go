//go:build v6

package web

import (
	"fmt"
	"strings"
)

// router 路由
type router struct {
	// trees 路由树（森林）
	// key 为 HTTP 方法
	// value 为 路由树
	trees map[string]*node
}

func newRouter() *router {
	return &router{trees: map[string]*node{}}
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
		r.trees[method] = root
	}

	// 特殊路由处理
	if path == "/" {
		if root.handler != nil {
			panic("route: 路由冲突[/]")
		}
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

// node 代表路由树的节点
// 路由树的匹配顺序是：
// 1. 静态完全匹配
// 2. 路径参数匹配：形式 :param_name
// 3. 通配符匹配：*
// 这是不回溯匹配
type node struct {
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

	// handler 业务逻辑
	handler HandleFunc
}

func (n *node) childOrCreate(path string) *node {
	// 通配符路由节点
	if path == "*" {
		if n.paramChild != nil {
			panic(fmt.Sprintf("route: 非法路由，已有参数路由，不能再注册通配符路由 [%s]", path))
		}
		if n.starChild == nil {
			n.starChild = &node{path: path}
		}
		return n.starChild
	}

	// 参数路径节点
	if path[0] == ':' {
		if n.starChild != nil {
			panic(fmt.Sprintf("route: 非法路由，已有通配符路由，不能再注册参数路由 [%s]", path))
		}
		if n.paramChild != nil {
			if n.paramChild.path != path {
				panic(fmt.Sprintf("route: 非法路由，参数路由冲突 [%s]", path))
			}
		} else {
			n.paramChild = &node{path: path}
		}
		return n.paramChild
	}

	// 如果 children 没有初始化，则进行初始化
	if n.children == nil {
		n.children = map[string]*node{}
	}
	child, ok := n.children[path]
	if !ok {
		// 没有的话需要创建一个 node
		child = &node{path: path}
		n.children[path] = child
	}
	return child
}

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
