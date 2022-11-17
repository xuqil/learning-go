//go:build v7

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
		// 静态路由测试用例
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
		// 通配符测试用例
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
		// 参数路由
		{
			method: http.MethodGet,
			path:   "/param/:id",
		},
		{
			method: http.MethodGet,
			path:   "/param/:id/detail",
		},
		{
			method: http.MethodGet,
			path:   "/param/:id/*",
		},
		// 正则路由
		{
			method: http.MethodDelete,
			path:   "/reg/:id(.*)",
		},
		{
			method: http.MethodDelete,
			path:   "/:name(^.+$)/abc",
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
								typ:     nodeTypeStatic,
							},
						},
						handler: mockHandler,
						typ:     nodeTypeStatic,
					},
					"order": {
						path: "order",
						children: map[string]*node{
							"detail": {
								path:    "detail",
								handler: mockHandler,
								typ:     nodeTypeStatic,
							},
						},
						starChild: &node{
							path:    "*",
							handler: mockHandler,
							typ:     nodeTypeAny,
						},
						typ: nodeTypeStatic,
					},
					"param": {
						path: "param",
						paramChild: &node{
							path: ":id",
							children: map[string]*node{
								"detail": {
									path:    "detail",
									handler: mockHandler,
									typ:     nodeTypeStatic,
								},
							},
							starChild: &node{
								path:    "*",
								handler: mockHandler,
								typ:     nodeTypeAny,
							},
							handler: mockHandler,
							typ:     nodeTypeParam,
						},
					},
				},
				starChild: &node{
					path: "*",
					children: map[string]*node{
						"abc": {
							path: "abc",
							starChild: &node{
								path:    "*",
								handler: mockHandler,
								typ:     nodeTypeAny,
							},
							handler: mockHandler,
							typ:     nodeTypeStatic,
						},
					},
					starChild: &node{
						path:    "*",
						handler: mockHandler,
						typ:     nodeTypeAny,
					},
					handler: mockHandler,
					typ:     nodeTypeAny,
				},
				handler: mockHandler,
				typ:     nodeTypeStatic,
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
								typ:     nodeTypeStatic,
							},
						},
					},
					"login": {
						path:    "login",
						handler: mockHandler,
						typ:     nodeTypeStatic,
					},
				},
				typ: nodeTypeStatic,
			},
			http.MethodDelete: {
				path: "/",
				children: map[string]*node{
					"reg": {
						path: "reg",
						typ:  nodeTypeStatic,
						regChild: &node{
							path:      ":id(.*)",
							paramName: "id",
							typ:       nodeTypeReg,
							handler:   mockHandler,
						},
					},
				},
				regChild: &node{
					path:      ":name(^.+$)",
					paramName: "name",
					typ:       nodeTypeReg,
					children: map[string]*node{
						"abc": {
							path:    "abc",
							handler: mockHandler,
						},
					},
				},
			},
		},
	}

	msg, ok := wantRouter.equal(&r)
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

	// 通配符节点重复注册
	r.addRoute(http.MethodGet, "/a/*/c", mockHandler)
	assert.PanicsWithValue(t, "route: 路由冲突[/a/*/c]", func() {
		r.addRoute(http.MethodGet, "/a/*/c", mockHandler)
	})
	// 多个 /
	assert.PanicsWithValue(t, "route: 非法路由，路由不能有连续的 /, [/a//b]", func() {
		r.addRoute(http.MethodGet, "/a//b", mockHandler)
	})
	assert.PanicsWithValue(t, "route: 非法路由，路由不能有连续的 /, [//a/b]", func() {
		r.addRoute(http.MethodGet, "//a/b", mockHandler)
	})
	// 同时注册通配符路由，参数路由，正则路由
	assert.PanicsWithValue(t, "route: 非法路由，已有通配符路由。不允许同时注册通配符路由和参数路由 [:id]", func() {
		r.addRoute(http.MethodGet, "/a/*", mockHandler)
		r.addRoute(http.MethodGet, "/a/:id", mockHandler)
	})
	r = newRouter()
	assert.PanicsWithValue(t, "route: 非法路由，已有通配符路由。不允许同时注册通配符路由和正则路由 [:id(.*)]", func() {
		r.addRoute(http.MethodGet, "/a/b/*", mockHandler)
		r.addRoute(http.MethodGet, "/a/b/:id(.*)", mockHandler)
	})
	r = newRouter()
	assert.PanicsWithValue(t, "route: 非法路由，已有通配符路由。不允许同时注册通配符路由和参数路由 [:id]", func() {
		r.addRoute(http.MethodGet, "/*", mockHandler)
		r.addRoute(http.MethodGet, "/:id", mockHandler)
	})
	r = newRouter()
	assert.PanicsWithValue(t, "route: 非法路由，已有路径参数路由。不允许同时注册通配符路由和参数路由 [*]", func() {
		r.addRoute(http.MethodGet, "/a/b/:id", mockHandler)
		r.addRoute(http.MethodGet, "/a/b/*", mockHandler)
	})
	r = newRouter()
	assert.PanicsWithValue(t, "route: 非法路由，已有路径参数路由。不允许同时注册通配符路由和参数路由 [*]", func() {
		r.addRoute(http.MethodGet, "/:id", mockHandler)
		r.addRoute(http.MethodGet, "/*", mockHandler)
	})
	r = newRouter()
	assert.PanicsWithValue(t, "route: 非法路由，已有路径参数路由。不允许同时注册正则路由和参数路由 [:id(.*)]", func() {
		r.addRoute(http.MethodGet, "/a/b/:id", mockHandler)
		r.addRoute(http.MethodGet, "/a/b/:id(.*)", mockHandler)
	})
	r = newRouter()

	assert.PanicsWithValue(t, "route: 非法路由，已有正则路由。不允许同时注册通配符路由和正则路由 [*]", func() {
		r.addRoute(http.MethodGet, "/a/b/:id(.*)", mockHandler)
		r.addRoute(http.MethodGet, "/a/b/*", mockHandler)
	})
	r = newRouter()
	assert.PanicsWithValue(t, "route: 非法路由，已有正则路由。不允许同时注册正则路由和参数路由 [:id]", func() {
		r.addRoute(http.MethodGet, "/a/b/:id(.*)", mockHandler)
		r.addRoute(http.MethodGet, "/a/b/:id", mockHandler)
	})
	// 参数冲突
	assert.PanicsWithValue(t, "route: 路由冲突，参数路由冲突，已有 :id，新注册 :name", func() {
		r.addRoute(http.MethodGet, "/a/b/c/:id", mockHandler)
		r.addRoute(http.MethodGet, "/a/b/c/:name", mockHandler)
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

	if n.typ != y.typ {
		return fmt.Sprintf("%s 节点类型不相等 x %d, y %d", n.path, n.typ, y.typ), false
	}

	if n.paramName != y.paramName {
		return fmt.Sprintf("%s 节点参数名字不相等 x %s, y %s", n.path, n.paramName, y.paramName), false
	}

	if len(n.children) != len(y.children) {
		return fmt.Sprintf("children 数量不相等, want: %d, now: %d", len(n.children), len(y.children)), false
	}

	if len(n.children) == 0 {
		return "", true
	}

	if n.paramChild != nil {
		str, ok := n.paramChild.equal(y.paramChild)
		if !ok {
			return fmt.Sprintf("%s 路径参数节点不匹配 %s", n.path, str), false
		}
	}

	if n.starChild != nil {
		str, ok := n.starChild.equal(y.starChild)
		if !ok {
			return fmt.Sprintf("%s 通配符节点不匹配 %s", n.path, str), false
		}
	}

	if n.regChild != nil {
		str, ok := n.regChild.equal(y.regChild)
		if !ok {
			return fmt.Sprintf("%s 路径参数节点不匹配 %s", n.path, str), false
		}
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

func Test_route_findRoute(t *testing.T) {
	testRoutes := []struct {
		method string
		path   string
	}{
		// 静态路由测试
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		// 通配符路由测试
		{
			method: http.MethodGet,
			path:   "/user/*/home",
		},
		{
			method: http.MethodPost,
			path:   "/order/*",
		},
		{
			method: http.MethodPost,
			path:   "/order/*/page/*",
		},
		//	参数路由
		{
			method: http.MethodGet,
			path:   "/param/:id",
		},
		{
			method: http.MethodGet,
			path:   "/param/:id/detail",
		},
		{
			method: http.MethodGet,
			path:   "/param/:id/*",
		},
	}

	var mockHandler HandleFunc = func(ctx *Context) {}
	testCases := []struct {
		name   string
		method string
		path   string

		wantFound bool
		wantMi    *matchInfo
	}{
		{
			name:   "method not found",
			method: http.MethodHead,
		},
		{
			name:   "path not found",
			method: http.MethodGet,
			path:   "/abc",
		},
		{
			name:      "root",
			method:    http.MethodGet,
			path:      "/",
			wantFound: true,
			wantMi: &matchInfo{
				n: &node{
					path:    "/",
					handler: mockHandler,
				},
			},
		},
		{
			name:      "user",
			method:    http.MethodGet,
			path:      "/user",
			wantFound: true,
			wantMi: &matchInfo{
				n: &node{
					path:    "user",
					handler: mockHandler,
				},
			},
		},
		{
			name:      "no handler",
			method:    http.MethodPost,
			path:      "/order",
			wantFound: true,
			wantMi: &matchInfo{
				n: &node{
					path: "order",
				},
			},
		},
		{
			name:      "two layer",
			method:    http.MethodPost,
			path:      "/order/create",
			wantFound: true,
			wantMi: &matchInfo{
				n: &node{
					path:    "create",
					handler: mockHandler,
				},
			},
		},
		// 通配符匹配
		{
			// 命中/order/*
			name:      "star match",
			method:    http.MethodPost,
			path:      "/order/delete",
			wantFound: true,
			wantMi: &matchInfo{
				n: &node{
					path:    "*",
					handler: mockHandler,
				},
			},
		},
		{
			// 命中通配符在中间的
			// /user/*/home
			name:      "star in middle",
			method:    http.MethodGet,
			path:      "/user/Tom/home",
			wantFound: true,
			wantMi: &matchInfo{
				n: &node{
					path:    "home",
					handler: mockHandler,
				},
			},
		},
		{
			// 比 /order/* 多了一段
			name:      "overflow",
			method:    http.MethodPost,
			path:      "/order/delete/123",
			wantFound: true,
			wantMi: &matchInfo{
				n: &node{
					path:    "*",
					handler: mockHandler,
				},
			},
		},
		{
			// 比 /order/*/page/* 多一段
			name:      "overflow2",
			method:    http.MethodPost,
			path:      "/order/123/page/10/delete",
			wantFound: true,
			wantMi: &matchInfo{
				n: &node{
					path:    "*",
					handler: mockHandler,
				},
			},
		},
		//	参数路由
		{
			name:      ":id",
			method:    http.MethodGet,
			path:      "/param/123",
			wantFound: true,
			wantMi: &matchInfo{
				n: &node{
					path:    ":id",
					handler: mockHandler,
				},
				pathParams: map[string]string{"id": "123"},
			},
		},
		{
			name:      ":id_detail",
			method:    http.MethodGet,
			path:      "/param/123/detail",
			wantFound: true,
			wantMi: &matchInfo{
				n: &node{
					path:    "detail",
					handler: mockHandler,
				},
				pathParams: map[string]string{"id": "123"},
			},
		},
		{
			name:      ":id*",
			method:    http.MethodGet,
			path:      "/param/111/user",
			wantFound: true,
			wantMi: &matchInfo{
				n: &node{
					path:    "*",
					handler: mockHandler,
				},
				pathParams: map[string]string{"id": "111"},
			},
		},
	}

	r := newRouter()
	for _, tr := range testRoutes {
		r.addRoute(tr.method, tr.path, mockHandler)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mi, found := r.findRoute(tc.method, tc.path)
			assert.Equal(t, tc.wantFound, found)
			if !found {
				return
			}
			assert.Equal(t, tc.wantMi.pathParams, mi.pathParams)
			n := mi.n
			assert.Equal(t, tc.wantMi.n.path, n.path)
			wantVal := reflect.ValueOf(tc.wantMi.n.handler)
			nVal := reflect.ValueOf(n.handler)
			assert.Equal(t, wantVal, nVal)
		})
	}
}
