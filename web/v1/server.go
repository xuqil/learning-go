package web

import (
	"net"
	"net/http"
)

type HandleFunc func(ctx Context)

// 确保一定实现了 Server 接口
var _ Server = &HTTPServer{}

// Server 核心 API 接口
type Server interface {
	http.Handler
	// Start 启动服务器
	// addr 是监听地址
	Start(addr string) error

	// AddRoute 路由注册功能
	// method 是 HTTP 方法
	// path 是路由
	// handleFunc 你的业务逻辑
	AddRoute(method string, path string, handleFunc HandleFunc)
	// AddRoute1 支持注册多个 handleFunc，没有必要提供
	//AddRoute1(method string, path string, handlerFunc ...HandleFunc)
}

//type HTTPSServer struct {
//	HTTPServer
//}

// HTTPServer 实现了 Server ，而 Server 由组合了 http.Handler
// 衍生 API
type HTTPServer struct {
}

// ServeHTTP 处理请求的入口
func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// 框架代码
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}
	h.serve(ctx)
}

func (h *HTTPServer) serve(ctx *Context) {
	//	接下来就是查看路由，并且执行命中的业务逻辑
}

func (h *HTTPServer) AddRoute(method string, path string, handlerFunc HandleFunc) {
	// 这里注册到路由树里面
	//panic("implement me")
}

func (h *HTTPServer) Get(path string, handleFunc HandleFunc) {
	h.AddRoute(http.MethodGet, path, handleFunc)
}

func (h *HTTPServer) POST(path string, handleFunc HandleFunc) {
	h.AddRoute(http.MethodPost, path, handleFunc)
}

func (h *HTTPServer) OPTIONS(path string, handleFunc HandleFunc) {
	h.AddRoute(http.MethodOptions, path, handleFunc)
}

func (h *HTTPServer) Start(addr string) error {
	// 可以进行生命周期管理
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	// 在这里，可以让用户注册所谓的 after start 回调
	// 比如往 admin 注册一下这个实例
	// 执行一些业务所需的前置条件
	return http.Serve(l, h)
}

//func (h *HTTPServer) Start1(addr string) error {
//	return http.ListenAndServe(addr, h)
//}
