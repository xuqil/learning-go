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

	h.POST("/form", func(ctx *Context) {
		_, err := ctx.Resp.Write([]byte(fmt.Sprintf("hello, %s", ctx.Req.URL.Path)))
		if err != nil {
			return
		}
	})

	h.GET("/values/:id", func(ctx *Context) {
		id, err := ctx.PathValue("id").AsInt64()
		if err != nil {
			ctx.Resp.WriteHeader(400)
			_, err := ctx.Resp.Write([]byte("id 输入不对"))
			if err != nil {
				return
			}
			return
		}
		_, err = ctx.Resp.Write([]byte(fmt.Sprintf("helo, %d", id)))
		if err != nil {
			return
		}
	})

	type User struct {
		Name string `json:"name"`
	}

	h.GET("/user/123", func(ctx *Context) {
		err := ctx.RespJSON(202, User{
			Name: "Tom",
		})
		if err != nil {
			return
		}
	})

	err := h.Start(":8081")
	if err != nil {
		return
	}
}
