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
func newRouter() *router {
	return &router{
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
		children := root.childrenOrCreate(seg)
		root = children
	}
	if root.handler != nil {
		panic(fmt.Sprintf("web: 路由冲突， 重复注册[%s]", path))
	}
	root.handler = handlerFunc
}

type node struct {
	path string

	// 子 path 到左节点的映射
	children map[string]*node

	//	业务逻辑
	handler HandleFunc
}

func (n *node) childrenOrCreate(seg string) *node {
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
