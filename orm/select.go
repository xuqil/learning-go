package orm

import (
	"context"
	"leanring-go/orm/internal/errs"
	"reflect"
	"strings"
)

type Selector[T any] struct {
	table string
	model *Model
	where []Predicate
	sb    *strings.Builder
	args  []any

	db *DB
}

func NewSelector[T any](db *DB) *Selector[T] {
	return &Selector[T]{
		sb: &strings.Builder{},
		db: db,
	}
}

func (s *Selector[T]) Build() (*Query, error) {
	var err error
	s.model, err = s.db.r.Get(new(T))
	if err != nil {
		return nil, err
	}
	sb := s.sb
	sb.WriteString("SELECT * FROM ")
	// 把表名拿到
	if s.table == "" {
		sb.WriteByte('`')
		sb.WriteString(s.model.tableName)
		sb.WriteByte('`')
	} else {
		//sb.WriteByte('`')
		sb.WriteString(s.table)
		//sb.WriteByte('`')
	}
	if len(s.where) > 0 {
		sb.WriteString(" WHERE ")
		p := s.where[0]
		for i := 1; i < len(s.where); i++ {
			p = p.And(s.where[i])
		}
		if err := s.buildExpression(p); err != nil {
			return nil, err
		}
	}

	sb.WriteByte(';')
	return &Query{
		SQL:  sb.String(),
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
		s.sb.WriteByte(' ')
		s.sb.WriteString(exp.op.String())
		s.sb.WriteByte(' ')
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
		fd, ok := s.model.fields[exp.name]
		// 字段（列）不对
		if !ok {
			return errs.NewErrUnknownField(exp.name)
		}
		s.sb.WriteByte('`')
		s.sb.WriteString(fd.colName)
		s.sb.WriteByte('`')
	case value:
		s.sb.WriteByte('?')
		s.addArg(exp.val)
	default:
		return errs.NewErrUnsupportedExpression(expr)
	}
	return nil
}
func (s *Selector[T]) addArg(val any) *Selector[T] {
	if s.args == nil {
		// 预估容量
		s.args = make([]any, 0, 8)
	}
	s.args = append(s.args, val)
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

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	q, err := s.Build()
	// 构造 SQL 失败
	if err != nil {
		return nil, err
	}

	db := s.db.db
	//	在这里，就是发起查询，并且处理结果集
	rows, err := db.QueryContext(ctx, q.SQL, q.Args...)
	// 查询错误
	if err != nil {
		return nil, err
	}

	// 确认有没有数据
	if !rows.Next() {
		// 返回 error，和 sql 包语义保持一致。sql.ErrNoRows
		return nil, ErrNoRows
	}

	// 获取 SELECT 出来了哪些列
	cs, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// 怎么利用 cs 解决类型问题和顺序问题
	tp := new(T)

	// 通过 cs 来构造 vals
	vals := make([]any, 0, len(cs))
	for _, c := range cs {
		// c 是列名
		for _, fd := range s.model.fields {
			if fd.colName == c {
				// 反射创建一个实例
				// 这里创建的实例时原本类型的指针类型
				// 例如 fd.Type = int，那么 val 是 *int
				val := reflect.New(fd.typ)
				vals = append(vals, val.Interface())
			}
		}
	}

	// 1、类型要匹配
	// 2、顺序要匹配
	rows.Scan(vals...)

	// 把 vals 塞进结果 tp 里面
	tpValue := reflect.ValueOf(tp)
	for i, c := range cs {
		for _, fd := range s.model.fields {
			if fd.colName == c {
				tpValue.Elem().FieldByName(fd.goName).
					Set(reflect.ValueOf(vals[i]).Elem())
			}
		}
	}

	// 在这里处理结果集
	return tp, err
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	q, err := s.Build()
	// 构造 SQL 失败
	if err != nil {
		return nil, err
	}

	db := s.db.db
	//	在这里，就是发起查询，并且处理结果集
	_, err = db.QueryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return nil, err
	}

	// 在这里处理结果集
	return nil, err

}
