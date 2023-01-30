package orm

import (
	"context"
	"database/sql"
)

type RawQuerier[T any] struct {
	core
	sess Session
	sql  string
	args []any
}

func (r RawQuerier[T]) Build() (*Query, error) {
	return &Query{
		SQL:  r.sql,
		Args: r.args,
	}, nil
}

func RawQuery[T any](sess Session, query string, args ...any) *RawQuerier[T] {
	c := sess.getCore()
	return &RawQuerier[T]{
		sql:  query,
		args: args,
		sess: sess,
		core: c,
	}
}

func (r RawQuerier[T]) Exec(ctx context.Context) Result {
	var err error
	r.model, err = r.r.Get(new(T))
	if err != nil {
		return Result{
			err: err,
		}
	}
	res := exec(ctx, r.sess, r.core, &QueryContext{
		Type:    "RAW",
		Builder: r,
		Model:   r.model,
	})
	var sqlRes sql.Result
	if res.Result != nil {
		sqlRes = res.Result.(sql.Result)
	}
	return Result{
		err: res.Err,
		res: sqlRes,
	}
}

func (r RawQuerier[T]) Get(ctx context.Context) (*T, error) {
	var err error
	r.model, err = r.r.Get(new(T))
	if err != nil {
		return nil, err
	}
	res := get[T](ctx, r.sess, r.core, &QueryContext{
		Type:    "RAW",
		Builder: r,
		Model:   r.model,
	})
	if res.Result != nil {
		return res.Result.(*T), res.Err
	}
	return nil, res.Err
}

func (r RawQuerier[T]) GetMulti(ctx context.Context) ([]*T, error) {
	//TODO implement me
	panic("implement me")
}