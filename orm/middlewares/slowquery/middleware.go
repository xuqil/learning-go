package querylog

import (
	"context"
	"leanring-go/orm"
	"log"
	"time"
)

type MiddlewareBuilder struct {
	// 慢查询阈值
	threshold time.Duration
	logFunc   func(query string, args []any)
}

func NewMiddlewareBuilder(threshold time.Duration) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		logFunc: func(query string, args []any) {
			log.Printf("SQL: %s, args: %v\n", query, args)
		},
		threshold: threshold,
	}
}

func (m *MiddlewareBuilder) LogFunc(fn func(query string, args []any)) *MiddlewareBuilder {
	m.logFunc = fn
	return m
}

func (m MiddlewareBuilder) Build() orm.Middleware {
	return func(next orm.Handler) orm.Handler {
		return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			startTime := time.Now()
			defer func() {
				duration := time.Since(startTime)
				// 不是慢查询
				if duration <= m.threshold {
					return
				}
				q, err := qc.Builder.Build()
				if err == nil {
					// 记录一下
					m.logFunc(q.SQL, q.Args)
				}
			}()

			return next(ctx, qc)
		}
	}
}
