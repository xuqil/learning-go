//go:build e2e

package accesslog

import (
	"fmt"
	"leanring-go/web"
	"testing"
)

func TestMiddlewareBuilder_BuildE2E(t *testing.T) {
	builder := MiddlewareBuilder{}
	fmt.Printf("%T\n", builder)
	mdl := builder.LogFunc(func(log string) {
		fmt.Println(log)
	}).Build()
	server := web.NewHTTPServer(web.ServerWithMiddleware(mdl))
	server.GET("/a/b/*", func(ctx *web.Context) {
		fmt.Println("hello, it's me")
		ctx.Resp.Write([]byte("hello, it's me"))
	})

	server.Start(":8081")
}
