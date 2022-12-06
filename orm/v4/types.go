package orm

import (
	"context"
	"database/sql"
)

// Querier 用于 SELECT 语句
type Querier[T any] interface {
	Get(ctx context.Context) (*T, error)
	GetMulti(ctx context.Context) ([]*T, error)

	// 或者
	//Get(ctx context.Context) (T, error)
	//GetMulti(ctx context.Context) ([]T, error)
}

// Executor 用于 INSERT, DELETE 和 UPDATE
type Executor interface {
	Exec(ctx context.Context) (sql.Result, error)
}

type QueryBuilder interface {
	Build() (*Query, error)

	// 或者
	//Build() (Query, error)
}

type Query struct {
	SQL  string
	Args []any
}

type TableName interface {
	TableName() string
}
