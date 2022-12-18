package orm

import (
	"leanring-go/orm/internal/errs"
	"reflect"
	"strings"
)

type Inserter[T any] struct {
	values []*T
	db     *DB
}

func NewInserter[T any](db *DB) *Inserter[T] {
	return &Inserter[T]{
		db: db,
	}
}

// Values 指定插入哪些数据
func (i *Inserter[T]) Values(vals ...*T) *Inserter[T] {
	i.values = vals
	return i
}

func (i *Inserter[T]) Build() (*Query, error) {
	if len(i.values) == 0 {
		return nil, errs.ErrInsertZeroRow
	}
	var sb strings.Builder
	sb.WriteString("INSERT INTO ")
	m, err := i.db.r.Get(i.values[0])
	if err != nil {
		return nil, err
	}
	// 拼接表名
	sb.WriteByte('`')
	sb.WriteString(m.TableName)
	sb.WriteByte('`')
	// 一定要显式指定列的顺序，不然不知道数据库中默认的顺序
	sb.WriteByte('(')
	// 不能通过遍历 map 类型（ FieldMap 和 ColMap ）获取 field 顺序，因为 map 类型是乱序的
	for idx, field := range m.Fields {
		if idx > 0 {
			sb.WriteByte(',')
		}
		sb.WriteByte('`')
		sb.WriteString(field.ColName)
		sb.WriteByte('`')
	}
	sb.WriteByte(')')

	// 拼接 Values
	sb.WriteString(" VALUES ")
	// 预估的参数数量：行数*列数
	args := make([]any, 0, len(i.values)*len(m.Fields))
	for j, val := range i.values {
		if j > 0 {
			sb.WriteByte(',')
		}
		sb.WriteByte('(')
		for idx, field := range m.Fields {
			if idx > 0 {
				sb.WriteByte(',')
			}
			sb.WriteByte('?')
			// 把参数读出来
			arg := reflect.ValueOf(val).Elem().FieldByName(field.GoName).Interface()
			args = append(args, arg)
		}
		sb.WriteByte(')')
	}
	sb.WriteByte(';')

	return &Query{
		SQL:  sb.String(),
		Args: args,
	}, nil
}
