package orm

import (
	"context"
	"leanring-go/orm/model"
)

type QueryContext struct {
	// 查询类型，标记增删改查
	Type string

	// 代表的是查询本身
	Builder QueryBuilder

	Model *model.Model
}

type QueryResult struct {
	//	Result 在不同的查询下类型是不同的
	Result any
	// 查询本身的错误
	Err error
}

type Handler func(ctx context.Context, qc *QueryContext) *QueryResult

type Middleware func(next Handler) Handler
