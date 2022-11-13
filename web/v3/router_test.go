package web

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func TestRouter_addRoute(t *testing.T) {
	// 第一步骤是构造路由树
	// 第第二个步骤是验证路由树
	testRouters := []struct {
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
			method: http.MethodGet,
			path:   "/order/detail/:id",
		},
		{
			method: http.MethodGet,
			path:   "/order/*",
		},
		{
			method: http.MethodGet,
			path:   "/*",
		},
		{
			method: http.MethodGet,
			path:   "/*/*",
		},
		{
			method: http.MethodGet,
			path:   "/*/abc",
		},
		{
			method: http.MethodGet,
			path:   "/*/abc/*",
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
	for _, route := range testRouters {
		r.addRoute(route.method, route.path, mockHandler)
	}

	//	在这里断言路由树和你的预期的一模一样
	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: {
				path:    "/",
				handler: mockHandler,
				children: map[string]*node{
					"user": {
						path:    "user",
						handler: mockHandler,
						children: map[string]*node{
							"home": {
								path:    "home",
								handler: mockHandler,
							},
						},
					},
					"order": {
						path: "order",
						children: map[string]*node{
							"detail": {
								path:    "detail",
								handler: mockHandler,
								paramChild: &node{
									path:    ":id",
									handler: mockHandler,
								},
							},
						},
						startChild: &node{
							path:    "*",
							handler: mockHandler,
						},
					},
				},
				startChild: &node{
					path:    "*",
					handler: mockHandler,
					children: map[string]*node{
						"abc": {
							path: "abc",
							startChild: &node{
								path:    "*",
								handler: mockHandler,
							},
							handler: mockHandler,
						},
					},
					startChild: &node{
						path:    "*",
						handler: mockHandler,
					},
				},
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
	// 断言两者相等（函数不可以比较，需要新定义方法帮助比较）
	msg, ok := wantRouter.equal(&r)
	//	msg, ok := r.equal(wantRouter)
	assert.True(t, ok, msg)

	r = newRouter()
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "", mockHandler)
	}, "web: 路径必须以 / 开头")
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/a/b/c/", mockHandler)
	}, "web: 路径不能以 / 结尾")
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/a//b/c/", mockHandler)
	}, "web: 路由不能有连续的 /")

	r = newRouter()
	r.addRoute(http.MethodGet, "/", mockHandler)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/", mockHandler)
	}, "路由冲突，重复注册[/]")

	r = newRouter()
	r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
	}, "路由冲突，重复注册[/a/b/c]")

	// 可用的 http method，要不要校验？AddRoute 改为 addRoute，变私有，用户不能使用
	// mockHandler 为 nil，要不要校验？用户决定，一般不会为 nil，如果为 nil 相当于没有注册路由

	r = newRouter()
	r.addRoute(http.MethodGet, "/a/*", mockHandler)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/a/:id", mockHandler)
	}, "web: 不允许同时注册路径参数和通配符匹配，已有通配符匹配")

	r = newRouter()
	r.addRoute(http.MethodGet, "/a/:id", mockHandler)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/a/*", mockHandler)
	}, "web: 不允许同时注册路径参数和通配符匹配，已有路径参数")
}

// string 返回错误信息，帮助排查
// bool 代表相等
func (r *router) equal(y *router) (string, bool) {
	for k, v := range r.trees {
		yv, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("目标 router 里面没有方法 %s 的路由树", k), false
		}
		if str, ok := v.equal(yv); !ok {
			return k + "-" + str, ok
		}
	}
	return "", true
}
func (n *node) equal(y *node) (string, bool) {
	if y == nil {
		return "目标节点为 nil", false
	}
	if n.path != y.path {
		return fmt.Sprintf("%s 节点 path 不相等 x %s, y %s", n.path, n.path, y.path), false
	}

	if n.startChild != nil {
		msg, ok := n.startChild.equal(y.startChild)
		if !ok {
			return msg, ok
		}
	}

	if n.paramChild != nil {
		msg, ok := n.paramChild.equal(y.paramChild)
		if !ok {
			return msg, ok
		}
	}

	// 比较两个方法是否相等
	nhv := reflect.ValueOf(n.handler)
	yhv := reflect.ValueOf(y.handler)
	if nhv != yhv {
		return fmt.Sprintf("%s 节点 handler 不相等 x %s, y %s", n.path, nhv.Type().String(), yhv.Type().String()), false
	}

	if len(n.children) != len(y.children) {
		return fmt.Sprintf("%s 子节点长度不等", n.path), false
	}
	if len(n.children) == 0 {
		return "", true
	}

	for k, v := range n.children {
		yv, ok := y.children[k]
		if !ok {
			return fmt.Sprintf("%s 目标节点缺少子节点 %s", n.path, k), false
		}
		if str, ok := v.equal(yv); !ok {
			return n.path + "-" + str, ok
		}
	}
	return "", true
}

func TestRouter_findRoute(t *testing.T) {
	testRouters := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodDelete,
			path:   "/",
		},
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
			method: http.MethodGet,
			path:   "/order/*",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodPost,
			path:   "/login",
		},
		{
			method: http.MethodPost,
			path:   "/login/:username",
		},
	}

	r := newRouter()
	var mockHandler HandleFunc = func(ctx *Context) {}
	for _, route := range testRouters {
		r.addRoute(route.method, route.path, mockHandler)
	}

	testCases := []struct {
		name string

		method string
		path   string

		wantFound bool
		info      *matchInfo
	}{
		{
			// 方法都不存在
			name:      "method not found",
			method:    http.MethodOptions,
			path:      "/order/detail",
			wantFound: false,
		},
		{
			// 完全命中
			name:      "order detail",
			method:    http.MethodGet,
			path:      "/order/detail",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					handler: mockHandler,
					path:    "detail",
				},
			},
		},
		{
			// 完全命中
			name:      "order start",
			method:    http.MethodGet,
			path:      "/order/abc",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					handler: mockHandler,
					path:    "*",
				},
			},
		},
		{
			// 命中了，但是没有 handler
			name:      "order",
			method:    http.MethodGet,
			path:      "/order",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					//handler: mockHandler,
					path: "order",
					children: map[string]*node{
						"detail": &node{
							path:    "detail",
							handler: mockHandler,
						},
					},
				},
			},
		},
		{
			// not found
			name:      "path not found",
			method:    http.MethodGet,
			path:      "/aaaaaa/bbbb",
			wantFound: false,
		},
		{
			// 根节点
			name:      "root",
			method:    http.MethodDelete,
			path:      "/",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					path:    "/",
					handler: mockHandler,
				},
			},
		},
		{
			// username 路径参数匹配
			name:      "login username",
			method:    http.MethodPost,
			path:      "/login/code",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					path:    ":username",
					handler: mockHandler,
				},
				pathParams: map[string]string{
					"username": "code",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			info, found := r.findRoute(tc.method, tc.path)
			assert.Equal(t, tc.wantFound, found)
			if !found {
				return
			}
			assert.Equal(t, tc.info.pathParams, info.pathParams)
			msg, ok := tc.info.n.equal(info.n)
			assert.True(t, ok, msg)
		})
	}
}
