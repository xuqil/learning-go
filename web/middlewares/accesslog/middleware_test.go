package accesslog

import (
	"fmt"
	"leanring-go/web"
	"net/http"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := MiddlewareBuilder{}
	mdl := builder.LogFunc(func(log string) {
		fmt.Println(log)
	}).Build()
	server := web.NewHTTPServer(web.ServerWithMiddleware(mdl))
	server.POST("/a/b/*", func(ctx *web.Context) {
		fmt.Println("hello, it's me")
	})

	req, err := http.NewRequest(http.MethodPost, "/a/b/c", nil)
	req.Host = "localhost"
	if err != nil {
		t.Fatal(err)
	}
	server.ServeHTTP(nil, req)
}
