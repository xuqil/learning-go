package recover

import "leanring-go/web"

type MiddlewareBuilder struct {
	StatusCode int
	ErrMsg     []byte
	LogFunc    func(ctx *web.Context)
}

func (m *MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			defer func() {
				if err := recover(); err != nil {
					ctx.RespStatusCode = m.StatusCode
					ctx.RespData = m.ErrMsg
					m.LogFunc(ctx)
				}
			}()
			next(ctx)
		}
	}
}
