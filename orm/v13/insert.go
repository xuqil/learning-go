package orm

import (
	"context"
	"database/sql"
	"leanring-go/orm/internal/errs"
	"leanring-go/orm/model"
)

type UpsertBuilder[T any] struct {
	i               *Inserter[T]
	conflictColumns []string
}

type Upsert struct {
	assigns         []Assignable
	conflictColumns []string
}

func (o *UpsertBuilder[T]) ConflictColumns(cols ...string) *UpsertBuilder[T] {
	o.conflictColumns = cols
	return o
}

func (o *UpsertBuilder[T]) Update(assigns ...Assignable) *Inserter[T] {
	o.i.onDuplicateKey = &Upsert{
		assigns:         assigns,
		conflictColumns: o.conflictColumns,
	}
	return o.i
}

// Assignable 标记接口
// 实现该接口意味着可以用于赋值语句
// 用于 UPDATE 和 UPSERT 中
type Assignable interface {
	assign()
}

type Inserter[T any] struct {
	builder
	sess    Session
	values  []*T
	columns []string

	//onConflict []Assignable
	onDuplicateKey *Upsert
}

func NewInserter[T any](sess Session) *Inserter[T] {
	c := sess.getCore()
	return &Inserter[T]{
		builder: builder{
			core:   c,
			quoter: c.dialect.quoter(),
		},
		sess: sess,
	}
}

// Values 指定插入哪些数据
func (i *Inserter[T]) Values(vals ...*T) *Inserter[T] {
	i.values = vals
	return i
}

func (i *Inserter[T]) Columns(cols ...string) *Inserter[T] {
	i.columns = cols
	return i
}

//func (i *Inserter[T]) Upsert(assigns ...Assignable) *Inserter[T] {
//	i.onConflict = assigns
//	return i
//}

func (i *Inserter[T]) OnDuplicateKey() *UpsertBuilder[T] {
	return &UpsertBuilder[T]{
		i: i,
	}
}

func (i *Inserter[T]) Build() (*Query, error) {
	if len(i.values) == 0 {
		return nil, errs.ErrInsertZeroRow
	}
	i.sb.WriteString("INSERT INTO ")
	if i.model == nil {
		m, err := i.r.Get(i.values[0])
		i.model = m
		if err != nil {
			return nil, err
		}
	}
	// 拼接表名
	i.quote(i.model.TableName)
	// 一定要显式指定列的顺序，不然不知道数据库中默认的顺序
	i.sb.WriteByte('(')
	fields := i.model.Fields
	// 用户指定了
	if len(i.columns) > 0 {
		fields = make([]*model.Field, 0, len(i.columns))
		for _, fd := range i.columns {
			fdMeta, ok := i.model.FieldMap[fd]
			if !ok {
				return nil, errs.NewErrUnknownField(fd)
			}
			fields = append(fields, fdMeta)
		}
	}

	// 不能通过遍历 map 类型（ FieldMap 和 ColMap ）获取 field 顺序，因为 map 类型是乱序的
	for idx, field := range fields {
		if idx > 0 {
			i.sb.WriteByte(',')
		}
		i.quote(field.ColName)
	}
	i.sb.WriteByte(')')

	// 拼接 Values
	i.sb.WriteString(" VALUES ")
	// 预估的参数数量：行数*列数
	i.args = make([]any, 0, len(i.values)*len(fields))
	for j, v := range i.values {
		if j > 0 {
			i.sb.WriteByte(',')
		}
		i.sb.WriteByte('(')
		val := i.creator(i.model, v)
		for idx, field := range fields {
			if idx > 0 {
				i.sb.WriteByte(',')
			}
			i.sb.WriteByte('?')
			// 把参数读出来
			var arg any
			arg, err := val.Field(field.GoName)
			if err != nil {
				return nil, err
			}
			i.addArg(arg)
		}
		i.sb.WriteByte(')')
	}

	if i.onDuplicateKey != nil {
		err := i.dialect.buildUpsert(&i.builder, i.onDuplicateKey)
		if err != nil {
			return nil, err
		}
	}

	i.sb.WriteByte(';')

	return &Query{SQL: i.sb.String(), Args: i.args}, nil
}

func (i *Inserter[T]) Exec(ctx context.Context) Result {
	var err error
	i.model, err = i.r.Get(new(T))
	if err != nil {
		return Result{
			err: err,
		}
	}
	root := i.execHandler
	for j := len(i.mdls) - 1; j >= 0; j-- {
		root = i.mdls[j](root)
	}
	res := root(ctx, &QueryContext{
		Type:    "INSERT",
		Builder: i,
		Model:   i.model,
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

var _ Handler = (&Inserter[int]{}).execHandler

func (i *Inserter[T]) execHandler(ctx context.Context, qc *QueryContext) *QueryResult {
	q, err := i.Build()
	if err != nil {
		return &QueryResult{
			Err: err,
			Result: Result{
				err: err,
			},
		}
	}
	res, err := i.sess.execContext(ctx, q.SQL, q.Args...)
	return &QueryResult{
		Err: err,
		Result: Result{
			err: err,
			res: res,
		},
	}
}
