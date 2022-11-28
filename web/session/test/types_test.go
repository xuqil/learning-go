package test

import (
	"leanring-go/web"
	"leanring-go/web/session"
	"leanring-go/web/session/cookie"
	"leanring-go/web/session/memory"
	"net/http"
	"testing"
	"time"
)

func TestSession(t *testing.T) {
	var m *session.Manager = &session.Manager{
		Propagator: cookie.NewPropagator(),
		Store:      memory.NewStore(time.Minute * 15),
		CtxSessKey: "sseKey",
	}
	// 简单的登录校验
	sever := web.NewHTTPServer(web.ServerWithMiddleware(func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			if ctx.Req.URL.Path == "/login" {
				// 用户准备登录
				next(ctx)
				return
			}
			_, err := m.GetSession(ctx)
			if err != nil {
				ctx.RespStatusCode = http.StatusUnauthorized
				ctx.RespData = []byte("请重新登录")
				return
			}

			// 刷新 session 的过期时间
			_ = m.RefreshSession(ctx)
			next(ctx)
		}
	}))

	// 登录
	sever.POST("/login", func(ctx *web.Context) {
		// 在这之前校验登录用户名和密码
		sess, err := m.InitSession(ctx)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("登陆失败")
			return
		}
		err = sess.Set(ctx.Req.Context(), "nickname", "cookie")
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("登陆失败")
			return
		}
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("登录成功")
		return
	})

	// 退出登录
	sever.POST("/logout", func(ctx *web.Context) {
		// 清理各种数据
		err := m.RemoveSession(ctx)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("退出失败")
			return
		}
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("退出成功")
	})

	sever.GET("/user", func(ctx *web.Context) {

		sess, _ := m.GetSession(ctx)

		// 把昵称从 session 里面拿出来
		val, _ := sess.Get(ctx.Req.Context(), "nickname")
		ctx.RespData = []byte(val.(string))
	})

	sever.Start(":8081")
}
