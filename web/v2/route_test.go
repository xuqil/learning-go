//go:build v2

package web

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func Test_route_addRoute(t *testing.T) {
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodPost,
			path:   "/login",
		},
	}

	var mockHandler HandleFunc = func(ctx *Context) {}
	r := newRouter()
	for _, tr := range testRoutes {
		r.addRoute(tr.method, tr.path, mockHandler)
	}

	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: {
				path: "/",
				children: map[string]*node{
					"user": {
						path: "user",
						children: map[string]*node{
							"home": {
								path:    "home",
								handler: mockHandler,
							},
						},
						handler: mockHandler,
					},
					"order": {
						path: "order",
						children: map[string]*node{
							"detail": {
								path:    "detail",
								handler: mockHandler,
							},
						},
					},
				},
				handler: mockHandler,
			},
			http.MethodPost: {
				path: "/",
				children: map[string]*node{
					"order": {
						path: "order",
						children: map[string]*node{
							"create": {
								path:    "create",
								handler: mockHandler,
							},
						},
					},
					"login": {
						path:    "login",
						handler: mockHandler,
					},
				},
			},
		},
	}

	msg, ok := wantRouter.equal(r)
	assert.True(t, ok, msg)

	// 非法用例
	r = newRouter()

	// 空字符串
	assert.PanicsWithValue(t, "route: 路由是空字符串", func() {
		r.addRoute(http.MethodGet, "", mockHandler)
	})

	// 前导没有 /
	assert.PanicsWithValue(t, "route: 路由必须以 / 开头", func() {
		r.addRoute(http.MethodGet, "a/b/c", mockHandler)
	})

	// 后缀有 /
	assert.PanicsWithValue(t, "route: 路由不能以 / 结尾", func() {
		r.addRoute(http.MethodGet, "/a/b/c/", mockHandler)
	})

	// 根节点重复注册
	r.addRoute(http.MethodGet, "/", mockHandler)
	assert.PanicsWithValue(t, "route: 路由冲突[/]", func() {
		r.addRoute(http.MethodGet, "/", mockHandler)
	})
	// 普通节点重复注册
	r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
	assert.PanicsWithValue(t, "route: 路由冲突[/a/b/c]", func() {
		r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
	})

	// 多个 /
	assert.PanicsWithValue(t, "route: 非法路由，路由不能有连续的 /, [/a//b]", func() {
		r.addRoute(http.MethodGet, "/a//b", mockHandler)
	})
	assert.PanicsWithValue(t, "route: 非法路由，路由不能有连续的 /, [//a/b]", func() {
		r.addRoute(http.MethodGet, "//a/b", mockHandler)
	})

}

// equal 判断两个 router 是否相等
func (r *router) equal(y *router) (string, bool) {
	if y == nil {
		return "路由不能为 nil", false
	}
	for k, v := range r.trees {
		n, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("method:%s 目标路由没有路由树 path:%s", k, v.path), false
		}
		msg, ok := v.equal(n)
		if !ok {
			return fmt.Sprintf("method:%s path:%s -> msg:%s", k, v.path, msg), false
		}
	}
	return "", true
}

// 判断两个 node 是否相等
func (n *node) equal(y *node) (string, bool) {
	if y == nil {
		return "目标节点为 nil", false
	}
	if n.path != y.path {
		return fmt.Sprintf("path 不相等, want: %s, now: %s", n.path, y.path), false
	}
	nhv := reflect.ValueOf(n.handler)
	yhv := reflect.ValueOf(y.handler)
	if nhv != yhv {
		return fmt.Sprintf("handler 不相等, want: %s, now: %s", nhv.Type().String(), yhv.Type().String()), false
	}
	if len(n.children) != len(y.children) {
		return fmt.Sprintf("children 数量不相等, want: %d, now: %d", len(n.children), len(y.children)), false
	}
	for k, v := range n.children {
		yn, ok := y.children[k]
		if !ok {
			return fmt.Sprintf("目标节点缺少 path, want: %s", k), false
		}
		msg, ok := v.equal(yn)
		if !ok {
			return fmt.Sprintf("目标节点不相等, want: %s, msg: %s", n.path, msg), ok
		}
	}
	return "", true
}
