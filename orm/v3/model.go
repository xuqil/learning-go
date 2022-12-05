package orm

import (
	"leanring-go/orm/internal/errs"
	"reflect"
	"unicode"
)

type model struct {
	tableName string
	fields    map[string]*field
}

type field struct {
	// 列名
	colName string
}

// registry 代表的是元数据的注册中心
type registry struct {
	models map[reflect.Type]*model
}

func newRegistry() *registry {
	return &registry{
		models: make(map[reflect.Type]*model, 64),
	}
}

func (r *registry) get(val any) (*model, error) {
	typ := reflect.TypeOf(val)
	m, ok := r.models[typ]
	if !ok {
		var err error
		m, err = r.parseModel(val)
		if err != nil {
			return nil, err
		}
		r.models[typ] = m
	}
	return m, nil
}

func (r *registry) parseModel(entity any) (*model, error) {
	typ := reflect.TypeOf(entity)
	// 限制只能用一级指针
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	typ = typ.Elem()
	numField := typ.NumField()
	fieldMap := make(map[string]*field, numField)
	for i := 0; i < numField; i++ {
		fd := typ.Field(i)
		fieldMap[fd.Name] = &field{
			colName: underscoreName(fd.Name),
		}
	}
	return &model{
		tableName: underscoreName(typ.Name()),
		fields:    fieldMap,
	}, nil
}

// underscoreName 驼峰转字符串命名
func underscoreName(tableName string) string {
	var buf []byte
	for i, v := range tableName {
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}

	}
	return string(buf)
}
