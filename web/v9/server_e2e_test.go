//go:build v9

package web

import (
	"fmt"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	h := NewHTTPServer()

	h.addRoute(http.MethodGet, "/user", func(ctx *Context) {
		fmt.Println("处理第一件事")
		fmt.Println("处理第二件事")
	})
	handler1 := func(ctx *Context) {
		fmt.Println("处理第一件事")
	}

	handler2 := func(ctx *Context) {
		fmt.Println("处理第二件事")
	}

	// 用户自己去管这种
	h.addRoute(http.MethodGet, "/user/detail", func(ctx *Context) {
		handler1(ctx)
		handler2(ctx)
	})

	h.GET("/login", func(ctx *Context) {

	})

	err := h.Start(":8081")
	if err != nil {
		return
	}
}
