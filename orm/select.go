package orm

import (
	"context"
	"leanring-go/orm/internal/errs"
	"strings"
)

type Selector[T any] struct {
	table string
	model *model
	where []Predicate
	sb    *strings.Builder
	args  []any
}

func (s *Selector[T]) Build() (*Query, error) {
	s.sb = &strings.Builder{}
	var err error
	s.model, err = parseModel(new(T))
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
	//TODO implement me
	panic("implement me")
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	//TODO implement me
	panic("implement me")
}
