package opentelemetry

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"leanring-go/web"
)

const instrumentationName = "leanring-go/web/middleware/opentelemetry"

type MiddlewareBuilder struct {
	Tracer trace.Tracer
}

//func NewMiddlewareBuilder(tracer trace.Tracer) *MiddlewareBuilder {
//	return &MiddlewareBuilder{
//		Tracer: tracer,
//	}
//}

func (m MiddlewareBuilder) Build() web.Middleware {
	if m.Tracer == nil {
		m.Tracer = otel.GetTracerProvider().Tracer(instrumentationName)
	}
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {

			reqCtx := ctx.Req.Context()
			// 尝试和客户端（上游）的 trace 结合在一起
			reqCtx = otel.GetTextMapPropagator().Extract(reqCtx, propagation.HeaderCarrier(ctx.Req.Header))

			reqCtx, span := m.Tracer.Start(reqCtx, "unknown")
			defer span.End()

			/*
				defer func() {
					// 执行完 next 才可以有值
					span.SetName(ctx.MatchRoute)

					// 把响应码加上去
					span.SetAttributes(attribute.Int("http.status", ctx.RespStatusCode))
					span.End()
				}()
			*/

			span.SetAttributes(attribute.String("http.method", ctx.Req.Method))
			span.SetAttributes(attribute.String("peer.hostname", ctx.Req.Host))
			span.SetAttributes(attribute.String("peer.address", ctx.Req.RemoteAddr))
			span.SetAttributes(attribute.String("http.url", ctx.Req.URL.String()))
			span.SetAttributes(attribute.String("http.schema", ctx.Req.URL.Scheme))
			span.SetAttributes(attribute.String("http.host", ctx.Req.Host))
			span.SetAttributes(attribute.String("http.proto", ctx.Req.Proto))
			span.SetAttributes(attribute.String("span.kind", "server"))
			span.SetAttributes(attribute.String("component", "web"))

			// 将 Req Context 关联 opentelemetry 的 Context
			ctx.Req = ctx.Req.WithContext(reqCtx)

			// 直接调用下一步
			next(ctx)

			// 执行完 next 才可以有值
			span.SetName(ctx.MatchRoute)

			// 把响应码加上去
			span.SetAttributes(attribute.Int("http.status", ctx.RespStatusCode))
		}
	}
}
