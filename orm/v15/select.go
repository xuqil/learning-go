package orm

import (
	"context"
	"leanring-go/orm/internal/errs"
)

// Selectable 是一个标记接口
// 它代表的是查找的列，或者聚合函数
type Selectable interface {
	selectable()
}

type Selector[T any] struct {
	builder
	table   string
	where   []Predicate
	columns []Selectable

	sess Session
}

func NewSelector[T any](sess Session) *Selector[T] {
	c := sess.getCore()
	return &Selector[T]{
		builder: builder{
			core:   c,
			quoter: c.dialect.quoter(),
		},
		sess: sess,
	}
}

func (s *Selector[T]) Build() (*Query, error) {
	if s.model == nil {
		var err error
		s.model, err = s.r.Get(new(T))
		if err != nil {
			return nil, err
		}
	}

	s.sb.WriteString("SELECT ")
	if err := s.buildColumns(); err != nil {
		return nil, err
	}
	s.sb.WriteString(" FROM ")
	// 把表名拿到
	if s.table == "" {
		s.sb.WriteByte('`')
		s.sb.WriteString(s.model.TableName)
		s.sb.WriteByte('`')
	} else {
		//r.sb.WriteByte('`')
		s.sb.WriteString(s.table)
		//r.sb.WriteByte('`')
	}
	if len(s.where) > 0 {
		s.sb.WriteString(" WHERE ")
		p := s.where[0]
		for i := 1; i < len(s.where); i++ {
			p = p.And(s.where[i])
		}
		if err := s.buildExpression(p); err != nil {
			return nil, err
		}
	}

	s.sb.WriteByte(';')
	return &Query{
		SQL:  s.sb.String(),
		Args: s.args,
	}, nil
}

func (s *Selector[T]) buildExpression(expr Expression) error {
	switch exp := expr.(type) {
	case nil:
	case Predicate:
		// 在这里处理 p
		// p.left 构建好
		// p.opt 构建好
		// p.right 构建好
		_, ok := exp.left.(Predicate)
		if ok {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(exp.left); err != nil {
			return err
		}
		if ok {
			s.sb.WriteByte(')')
		}
		if exp.op != "" {
			s.sb.WriteByte(' ')
			s.sb.WriteString(exp.op.String())
			s.sb.WriteByte(' ')
		}
		_, ok = exp.right.(Predicate)
		if ok {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(exp.right); err != nil {
			return err
		}
		if ok {
			s.sb.WriteByte(')')
		}
	case Column:
		exp.alias = ""
		return s.buildColumn(exp)
	case value:
		s.sb.WriteByte('?')
		s.addArg(exp.val)
	case RawExpr:
		s.sb.WriteByte('(')
		s.sb.WriteString(exp.raw)
		s.addArg(exp.args...)
		s.sb.WriteByte(')')
	default:
		return errs.NewErrUnsupportedExpression(expr)
	}
	return nil
}

func (s *Selector[T]) buildColumns() error {
	if len(s.columns) == 0 {
		// 没有指定列
		s.sb.WriteByte('*')
		return nil
	}

	for i, col := range s.columns {
		if i > 0 {
			s.sb.WriteByte(',')
		}
		switch c := col.(type) {
		case Column:
			err := s.buildColumn(c)
			if err != nil {
				return err
			}
		case Aggregate:
			// 聚合函数名
			s.sb.WriteString(c.fn)
			s.sb.WriteByte('(')
			err := s.buildColumn(Column{name: c.arg})
			if err != nil {
				return err
			}
			s.sb.WriteByte(')')
			// 呼和函数本身的别名
			if c.alias != "" {
				s.sb.WriteString(" AS `")
				s.sb.WriteString(c.alias)
				s.sb.WriteByte('`')
			}
		case RawExpr:
			s.sb.WriteString(c.raw)
			s.addArg(c.args...)
		}
	}

	return nil
}
func (s *Selector[T]) buildColumn(c Column) error {
	fd, ok := s.model.FieldMap[c.name]
	// 字段（列）不对
	if !ok {
		return errs.NewErrUnknownField(c.name)
	}
	s.sb.WriteByte('`')
	s.sb.WriteString(fd.ColName)
	s.sb.WriteByte('`')
	if c.alias != "" {
		s.sb.WriteString(" AS `")
		s.sb.WriteString(c.alias)
		s.sb.WriteByte('`')
	}
	return nil
}

func (s *Selector[T]) Select(cols ...Selectable) *Selector[T] {
	s.columns = cols
	return s
}

func (s *Selector[T]) From(table string) *Selector[T] {
	s.table = table
	return s
}

func (s *Selector[T]) Where(ps ...Predicate) *Selector[T] {
	s.where = ps
	return s
}

//func (s *Selector[T]) GetV1(ctx context.Context) (*T, error) {
//	q, err := s.Build()
//	// 构造 SQL 失败
//	if err != nil {
//		return nil, err
//	}
//
//	db := s.db.db
//	//	在这里，就是发起查询，并且处理结果集
//	rows, err := db.QueryContext(ctx, q.SQL, q.Args...)
//	// 查询错误
//	if err != nil {
//		return nil, err
//	}
//
//	// 确认有没有数据
//	if !rows.Next() {
//		// 返回 error，和 sql 包语义保持一致。sql.ErrNoRows
//		return nil, ErrNoRows
//	}
//
//	// 获取 SELECT 出来了哪些列
//	cs, err := rows.Columns()
//	if err != nil {
//		return nil, err
//	}
//
//	var vals []any
//	tp := new(T)
//	// 起始地址
//	address := reflect.ValueOf(tp).UnsafePointer()
//	for _, c := range cs {
//		// c 是列名
//		fd, ok := s.model.ColumnMap[c]
//		if !ok {
//			return nil, errs.NewErrUnknownColumn(c)
//		}
//
//		// 计算字段的地址
//		// 起始地址 + 偏移量
//		fdAddress := unsafe.Pointer(uintptr(address) + fd.Offset)
//
//		// 反射在特定的地址上，创建一个特定类型的实例
//		// 这里创建的实例时原本类型的指针类型
//		// 例如 fd.Type = int，那么 val 是 *int
//		val := reflect.NewAt(fd.Typ, fdAddress)
//		vals = append(vals, val.Interface())
//	}
//
//	err = rows.Scan(vals...)
//	return tp, err
//}

//func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
//	var err error
//	s.model, err = s.r.Get(new(T))
//	if err != nil {
//		return nil, err
//	}
//
//	root := s.getHandler
//	for i := len(s.mdls) - 1; i >= 0; i-- {
//		root = s.mdls[i](root)
//	}
//	res := root(ctx, &QueryContext{
//		Type:    "SELECT",
//		Builder: s,
//		Model:   s.model,
//	})
//	if res.Result != nil {
//		return res.Result.(*T), res.Err
//	}
//	return nil, res.Err
//}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	var err error
	s.model, err = s.r.Get(new(T))
	if err != nil {
		return nil, err
	}
	res := get[T](ctx, s.sess, s.core, &QueryContext{
		Type:    "SELECT",
		Builder: s,
		Model:   s.model,
	})
	if res.Result != nil {
		return res.Result.(*T), res.Err
	}
	return nil, res.Err
}

//func getHandler[T any](ctx context.Context, sess Session, c core, qc *QueryContext) *QueryResult {
//	q, err := qc.Builder.Build()
//	// 构造 SQL 失败
//	if err != nil {
//		return &QueryResult{
//			Err: err,
//		}
//	}
//
//	//	在这里，就是发起查询，并且处理结果集
//	rows, err := sess.queryContext(ctx, q.SQL, q.Args...)
//	// 查询错误
//	if err != nil {
//		return &QueryResult{
//			Err: err,
//		}
//	}
//
//	// 确认有没有数据
//	if !rows.Next() {
//		// 返回 error，和 sql 包语义保持一致。sql.ErrNoRows
//		return &QueryResult{
//			Err: ErrNoRows,
//		}
//	}
//
//	tp := new(T)
//	val := c.creator(c.model, tp)
//	err = val.SetColumns(rows)
//
//	// 接口定义好之后：法1.用新接口的方法改造上层； 法2.提供不同的实现
//	return &QueryResult{
//		Err:    err,
//		Result: tp,
//	}
//}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	q, err := s.Build()
	// 构造 SQL 失败
	if err != nil {
		return nil, err
	}

	//	在这里，就是发起查询，并且处理结果集
	_, err = s.sess.queryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return nil, err
	}

	// 在这里处理结果集
	return nil, err

}