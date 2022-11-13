package web

import (
	"fmt"
	"strings"
)

// router 依赖支持对路由树的操作
// 代表路由树（森林）
type router struct {
	// http method => 路由树根节点
	trees map[string]*node
}

// newRouter 初始化 router
func newRouter() router {
	return router{
		trees: map[string]*node{},
	}
}

// AddRoute 注册路由
// method 是方法
// path 必须以 / 开始并且结尾不能有 /，中间也不允许有连续的 /
func (r *router) addRoute(method string, path string, handlerFunc HandleFunc) {
	// 这里注册到路由树里面
	// 开头不能没有 /
	if path == "" {
		panic("web: 路径不能为空字符串")
	}
	// 结尾
	if path[0] != '/' {
		panic("web: 路由必须以 / 开头")
	}

	if path != "/" && path[len(path)-1] == '/' {
		panic("web: 路由不能以 / 结尾")
	}

	// 中间连续 //，可以用 strings.contains("//")检查

	root, ok := r.trees[method]
	if !ok {
		//	说明没有根节点，需要创建根节点
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}
	// 根节点特殊处理
	if path == "/" {
		// 跟节点重复注册
		if root.handler != nil {
			panic("web: 路由冲突，重复注册[/]")
		}
		root.handler = handlerFunc
		return
	}
	//	切割 path，去掉前缀“/”：path[1:]
	for _, seg := range strings.Split(path[1:], "/") {
		if seg == "" {
			panic("web: 路由不能有连续的 x")
		}
		child := root.childrenOrCreate(seg)
		root = child
	}
	if root.handler != nil {
		panic(fmt.Sprintf("web: 路由冲突， 重复注册[%s]", path))
	}
	root.handler = handlerFunc
}

type node struct {
	path string

	// 静态匹配的节点
	// 子 path 到左节点的映射
	children map[string]*node

	// 通配符匹配的节点
	startChild *node

	// 路径参数
	paramChild *node

	//	业务逻辑
	handler HandleFunc
}

func (r *router) findRoute(method string, path string) (*node, bool) {
	// 沿着树深度优先搜索
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	// 根节点特殊处理
	if path == "/" {
		return root, true
	}

	// 去除前置后置 /
	path = strings.Trim(path, "/")
	for _, seg := range strings.Split(path, "/") {
		child, found := root.childOf(seg)
		if !found {
			return nil, false
		}
		root = child
	}
	//return root, root.handler != nil
	// 确实有这个节点，但不能确定有 handler
	return root, true
}

func (n *node) childrenOrCreate(seg string) *node {
	if seg[0] == ':' {
		if n.startChild != nil {
			panic("web: 不允许同时注册路径参数和通配符匹配，已有通配符匹配")
		}
		n.paramChild = &node{
			path: seg,
		}
		return n.paramChild
	}

	if seg == "*" {
		if n.paramChild != nil {
			panic("web: 不允许同时注册路径参数和通配符匹配，已有路径参数")
		}
		if n.startChild == nil {
			n.startChild = &node{path: "*"}
		}
		return n.startChild
	}
	if n.children == nil {
		n.children = map[string]*node{}
	}
	res, ok := n.children[seg]
	if !ok {
		res = &node{
			path: seg,
		}
		n.children[seg] = res
	}
	return res
}

// childOf 优先考虑静态匹配，匹配不上，再考虑通配符
func (n *node) childOf(path string) (*node, bool) {
	if n.children == nil {
		if n.paramChild != nil {
			return n.paramChild, true
		}
		return n.startChild, n.startChild != nil
	}
	child, ok := n.children[path]
	if !ok {
		if n.paramChild != nil {
			return n.paramChild, true
		}
		return n.startChild, n.startChild != nil
	}
	return child, ok
}
