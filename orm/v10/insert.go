package orm

import (
	"leanring-go/orm/internal/errs"
	"leanring-go/orm/model"
	"reflect"
)

type OnDuplicateKeyBuilder[T any] struct {
	i *Inserter[T]
}

type OnDuplicateKey struct {
	assigns []Assignable
}

func (o *OnDuplicateKeyBuilder[T]) Update(assigns ...Assignable) *Inserter[T] {
	o.i.onDuplicateKey = &OnDuplicateKey{
		assigns: assigns,
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
	values  []*T
	db      *DB
	columns []string

	//onConflict []Assignable
	onDuplicateKey *OnDuplicateKey
}

func NewInserter[T any](db *DB) *Inserter[T] {
	return &Inserter[T]{
		builder: builder{
			dialect: db.dialect,
			quoter:  db.dialect.quoter(),
		},
		db: db,
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

//func (i *Inserter[T]) OnDuplicateKey(assigns ...Assignable) *Inserter[T] {
//	i.onConflict = assigns
//	return i
//}

func (i *Inserter[T]) OnDuplicateKey() *OnDuplicateKeyBuilder[T] {
	return &OnDuplicateKeyBuilder[T]{
		i: i,
	}
}

func (i *Inserter[T]) Build() (*Query, error) {
	if len(i.values) == 0 {
		return nil, errs.ErrInsertZeroRow
	}
	i.sb.WriteString("INSERT INTO ")
	m, err := i.db.r.Get(i.values[0])
	i.model = m
	if err != nil {
		return nil, err
	}
	// 拼接表名
	i.quote(m.TableName)
	// 一定要显式指定列的顺序，不然不知道数据库中默认的顺序
	i.sb.WriteByte('(')
	fields := m.Fields
	// 用户指定了
	if len(i.columns) > 0 {
		fields = make([]*model.Field, 0, len(i.columns))
		for _, fd := range i.columns {
			fdMeta, ok := m.FieldMap[fd]
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
	for j, val := range i.values {
		if j > 0 {
			i.sb.WriteByte(',')
		}
		i.sb.WriteByte('(')
		for idx, field := range fields {
			if idx > 0 {
				i.sb.WriteByte(',')
			}
			i.sb.WriteByte('?')
			// 把参数读出来
			arg := reflect.ValueOf(val).Elem().FieldByName(field.GoName).Interface()
			i.addArg(arg)
		}
		i.sb.WriteByte(')')
	}

	if i.onDuplicateKey != nil {
		err = i.dialect.buildOnDuplicateKey(&i.builder, i.onDuplicateKey)
		if err != nil {
			return nil, err
		}
	}

	i.sb.WriteByte(';')

	return &Query{SQL: i.sb.String(), Args: i.args}, nil
}
