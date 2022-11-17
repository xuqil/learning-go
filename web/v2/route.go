//go:build v2

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

// node 路由树
type node struct {
	// path 路由
	path string

	// children 子节点
	// key 为 path
	// value 是 子节点
	children map[string]*node

	// handler 业务逻辑
	handler HandleFunc
}

func (n *node) childOrCreate(path string) *node {
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
