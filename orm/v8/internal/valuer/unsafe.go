package valuer

import (
	"database/sql"
	"leanring-go/orm/internal/errs"
	"leanring-go/orm/model"
	"reflect"
	"unsafe"
)

type unsafeValue struct {
	model *model.Model
	// 起始地址
	address unsafe.Pointer
}

var _ Creator = NewUnsafeValue

func NewUnsafeValue(model *model.Model, val any) Value {
	address := reflect.ValueOf(val).UnsafePointer()
	return unsafeValue{
		model:   model,
		address: address,
	}
}

func (r unsafeValue) SetColumns(rows *sql.Rows) error {
	// 获取 SELECT 出来了哪些列
	cs, err := rows.Columns()
	if err != nil {
		return err
	}

	var vals []any
	// 起始地址

	for _, c := range cs {
		// c 是列名
		fd, ok := r.model.ColumnMap[c]
		if !ok {
			return errs.NewErrUnknownColumn(c)
		}

		// 计算字段的地址
		// 起始地址 + 偏移量
		fdAddress := unsafe.Pointer(uintptr(r.address) + fd.Offset)

		// 反射在特定的地址上，创建一个特定类型的实例
		// 这里创建的实例时原本类型的指针类型
		// 例如 fd.Type = int，那么 val 是 *int
		val := reflect.NewAt(fd.Typ, fdAddress)
		vals = append(vals, val.Interface())
	}

	err = rows.Scan(vals...)
	return err
}
