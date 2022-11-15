//go:build v4

package web

import (
	"net"
	"net/http"
)

type HandleFunc func(ctx *Context)

// Server 框架核心接口
type Server interface {
	http.Handler

	// Start 启动服务
	Start(address string) error

	// AddRoute 路由注册功能
	// method 是 HTTP 方法
	// path 是路由
	// handleFunc 是业务逻辑
	addRoute(method string, path string, handleFunc HandleFunc)
}

// 确保 HTTPServer 一定实现了 Server
var _ Server = &HTTPServer{}

type HTTPServer struct {
	router
}

// ServeHTTP 服务请求入口
func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}
	h.serve(ctx)
}

// serve 路由匹配和业务逻辑执行
func (h *HTTPServer) serve(ctx *Context) {

}

func (h *HTTPServer) GET(path string, handler HandleFunc) {
	h.addRoute(http.MethodGet, path, handler)
}

func (h *HTTPServer) POST(path string, handler HandleFunc) {
	h.addRoute(http.MethodPost, path, handler)
}

func (h *HTTPServer) OPTIONS(path string, handler HandleFunc) {
	h.addRoute(http.MethodOptions, path, handler)
}

func (h *HTTPServer) DELETE(path string, handler HandleFunc) {
	h.addRoute(http.MethodDelete, path, handler)
}

// Start 由框架管理 HTTP Server
func (h *HTTPServer) Start(address string) error {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	return http.Serve(l, h)
}
