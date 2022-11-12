package web

import "strings"

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
func (r *router) AddRoute(method string, path string, handlerFunc HandleFunc) {
	// 这里注册到路由树里面
	root, ok := r.trees[method]
	if !ok {
		//	说明没有根节点，需要创建根节点
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}
	//	切割 path，去掉前缀“/”：path[1:]
	for _, seg := range strings.Split(path[1:], "/") {
		children := root.childrenOrCreate(seg)
		root = children
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
