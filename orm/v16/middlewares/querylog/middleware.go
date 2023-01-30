package querylog

import (
	"context"
	"leanring-go/orm"
	"log"
)

type MiddlewareBuilder struct {
	logFunc func(query string, args []any)
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		logFunc: func(query string, args []any) {
			log.Printf("SQL: %s, args: %v\n", query, args)
		},
	}
}

func (m *MiddlewareBuilder) LogFunc(fn func(query string, args []any)) *MiddlewareBuilder {
	m.logFunc = fn
	return m
}

func (m MiddlewareBuilder) Build() orm.Middleware {
	return func(next orm.Handler) orm.Handler {
		return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			q, err := qc.Builder.Build()
			if err != nil {
				return &orm.QueryResult{
					Err: err,
				}
			}

			m.logFunc(q.SQL, q.Args)
			return next(ctx, qc)
		}
	}
}