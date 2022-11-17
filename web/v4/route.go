//go:build v4

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

func newRouter() router {
	return router{trees: map[string]*node{}}
}

// addRoute 路由注册功能
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
func (r *router) findRoute(method string, path string) (*node, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	// 特殊路径
	if path == "/" {
		return root, true
	}
	// 切割 path 并匹配路由
	for _, seg := range strings.Split(strings.Trim(path, "/"), "/") {
		root, ok = root.childOf(seg)
		if !ok {
			return nil, ok
		}
	}
	return root, true
}

// node 代表路由树的节点
// 路由树的匹配顺序是：
// 1. 静态完全匹配
// 2. 通配符匹配
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

	// handler 业务逻辑
	handler HandleFunc
}

func (n *node) childOrCreate(path string) *node {
	// 通配符路由节点
	if path == "*" {
		if n.starChild == nil {
			n.starChild = &node{path: path}
		}
		return n.starChild
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

func (n *node) childOf(path string) (*node, bool) {
	if n.children == nil {
		return n.starChild, n.starChild != nil
	}
	child, ok := n.children[path]
	if !ok {
		return n.starChild, n.starChild != nil
	}
	return child, ok
}
