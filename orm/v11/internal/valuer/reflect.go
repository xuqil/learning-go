package valuer

import (
	"database/sql"
	"leanring-go/orm/internal/errs"
	"leanring-go/orm/model"
	"reflect"
)

type reflectValue struct {
	model *model.Model
	// 对应于 T 的指针
	val reflect.Value
}

var _ Creator = NewReflectValue

func NewReflectValue(model *model.Model, val any) Value {
	return reflectValue{
		model: model,
		val:   reflect.ValueOf(val).Elem(),
	}
}

func (r reflectValue) Field(name string) (any, error) {
	val := r.val.FieldByName(name)
	return val.Interface(), nil
}

func (r reflectValue) SetColumns(rows *sql.Rows) error {

	// 获取 SELECT 出来了哪些列
	cs, err := rows.Columns()
	if err != nil {
		return err
	}

	// 怎么利用 cs 解决类型问题和顺序问题
	// 通过 cs 来构造 vals
	vals := make([]any, 0, len(cs))
	valElem := make([]reflect.Value, 0, len(cs))
	for _, c := range cs {
		// c 是列名
		fd, ok := r.model.ColumnMap[c]
		if !ok {
			return errs.NewErrUnknownColumn(c)
		}

		// 反射创建一个实例
		// 这里创建的实例时原本类型的指针类型
		// 例如 fd.Type = int，那么 val 是 *int
		val := reflect.New(fd.Typ)
		vals = append(vals, val.Interface())
		// 记得要调用 Elem, 因为 fd.Type = int，那么 val 是 *int
		valElem = append(valElem, val.Elem())
	}

	// 1、类型要匹配
	// 2、顺序要匹配
	err = rows.Scan(vals...)
	if err != nil {
		return err
	}

	// 把 vals 塞进结果 tp 里面
	tpValueElem := r.val
	for i, c := range cs {
		// c 是列名
		fd, ok := r.model.ColumnMap[c]
		if !ok {
			return errs.NewErrUnknownColumn(c)
		}
		//tpValue.Elem().FieldByName(fd.goName).
		//	Set(reflect.ValueOf(vals[i]).Elem())
		tpValueElem.FieldByName(fd.GoName).Set(valElem[i])
	}

	// 在这里处理结果集
	return nil
}
