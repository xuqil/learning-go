package web

import (
	"fmt"
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

	mdls []Middleware

	log func(msg string, args ...any)

	tplEngine TemplateEngin
}

type HTTPServerOption func(server *HTTPServer)

func NewHTTPServer(opts ...HTTPServerOption) *HTTPServer {
	res := &HTTPServer{
		router: newRouter(),
		log: func(msg string, args ...any) {
			fmt.Printf(msg, args...)
		},
	}
	for _, opt := range opts {
		opt(res)
	}

	return res
}

func ServerWithTemplateEngine(tplEngine TemplateEngin) HTTPServerOption {
	return func(server *HTTPServer) {
		server.tplEngine = tplEngine
	}
}

func ServerWithMiddleware(mdls ...Middleware) HTTPServerOption {
	return func(server *HTTPServer) {
		server.mdls = mdls
	}
}

// ServeHTTP 服务请求入口
func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:       request,
		Resp:      writer,
		tplEngine: h.tplEngine,
	}
	// 最后一个是这个
	root := h.serve

	// 然后这里就是利用最后一个不断往前回溯组装链条
	// 从后往前，把后一个作为前一个的 next 构造好链条
	for i := len(h.mdls) - 1; i >= 0; i-- {
		root = h.mdls[i](root)
	}
	// 把 RespData 和 RespStatusCode 刷新到响应里面
	var m Middleware = func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			next(ctx)
			h.flashResp(ctx)
		}
	}
	root = m(root)
	//h.serve(ctx)
	// 这里执行的时候，就是从前往后了
	root(ctx)
}

func (h *HTTPServer) flashResp(ctx *Context) {
	if ctx.RespStatusCode != 0 {
		ctx.Resp.WriteHeader(ctx.RespStatusCode)
	}
	n, err := ctx.Resp.Write(ctx.RespData)
	if err != nil || n != len(ctx.RespData) {
		h.log("写入响应失败 %v", err)
	}
}

// serve 路由匹配和业务逻辑执行
func (h *HTTPServer) serve(ctx *Context) {
	mi, ok := h.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || mi.n == nil || mi.n.handler == nil {
		//ctx.Resp.WriteHeader(404)
		//_, err := ctx.Resp.Write([]byte("Not Found"))
		//if err != nil {
		//	return
		//}

		ctx.RespStatusCode = 404
		ctx.RespData = []byte("Not Found")
		return
	}
	ctx.PathParams = mi.pathParams
	ctx.MatchRoute = mi.n.route
	mi.n.handler(ctx)
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
