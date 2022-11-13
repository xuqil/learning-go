//go:build e2e

package web

import (
	"fmt"
	"testing"
)

// func TestServer(t *testing.T) {
//
//		//  1. 方法1，完全委托给 Http 包管理。由 Http 管理的话控制力比较差
//		// 这个 Handler 就是我们跟 Http 包的结构体
//		var h Server // NewServer
//		//var h = &HTTPServer{}
//		err := http.ListenAndServe(":8081", h)
//		if err != nil {
//			return
//		}
//		err = http.ListenAndServeTLS(":443", "", "", h)
//		if err != nil {
//			return
//		}
//
//		//	2. 方法2，自己手动管
//		h.Start("8082")
//	}
func TestServer(t *testing.T) {
	var h = NewHTTPServer()
	/*h.addRoute(http.MethodGet, "/user", func(ctx *Context) {
		fmt.Println("处理一件事")
	})

	handler1 := func(ctx *Context) {
		fmt.Println("处理第一件事")
	}
	handler2 := func(ctx *Context) {
		fmt.Println("处理第二件事")
	}
	// 允许用户注册多个 handleFunc
	// 用户自己管理多个 handleFunc
	h.addRoute(http.MethodGet, "/user", func(ctx *Context) {
		handler1(ctx)
		handler2(ctx)
	})*/

	h.Get("/hello", func(ctx *Context) {
		fmt.Println("GET")
	})

	h.Get("/user", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, user"))
	})

	h.Get("/order/detail", func(ctx *Context) {
		ctx.Resp.Write([]byte(fmt.Sprintf("hello, order detail")))
	})

	h.Get("/order/abc", func(ctx *Context) {
		ctx.Resp.Write([]byte(fmt.Sprintf("hello,%s", ctx.Req.URL.Path)))
	})

	//	2. 方法2，自己手动管
	err := h.Start(":8081")
	if err != nil {
		return
	}
}
