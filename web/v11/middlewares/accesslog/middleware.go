package accesslog

import (
	"encoding/json"
	"leanring-go/web"
)

type MiddlewareBuilder struct {
	logFunc func(log string)
}

func (m *MiddlewareBuilder) LogFunc(fn func(log string)) *MiddlewareBuilder {
	m.logFunc = fn
	return m
}

func (m *MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			// 要记录请求
			defer func() {
				l := accessLog{
					Host:       ctx.Req.Host,
					Route:      ctx.MatchRoute,
					HttpMethod: ctx.Req.Method,
					Path:       ctx.Req.URL.Path,
				}
				data, _ := json.Marshal(l)
				m.logFunc(string(data))
			}()
			next(ctx)
		}
	}
}

type accessLog struct {
	Host string `json:"host,omitempty"`
	// 命中的路由
	Route      string `json:"route,omitempty"`
	HttpMethod string `json:"http_method,omitempty"`
	Path       string `json:"path,omitempty"`
}
