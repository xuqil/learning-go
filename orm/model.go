package orm

import (
	"leanring-go/orm/internal/errs"
	"reflect"
	"sync"
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
	// 读写锁
	//lock   sync.RWMutex
	models sync.Map
}

func newRegistry() *registry {
	return &registry{}
}

// get 会存在覆盖问题
func (r *registry) get(val any) (*model, error) {
	typ := reflect.TypeOf(val)
	m, ok := r.models.Load(typ)
	if ok {
		return m.(*model), nil
	}
	m, err := r.parseModel(val)
	if err != nil {
		return nil, err
	}
	r.models.Store(typ, m)
	return m.(*model), nil
}

// get 使用 double read 读写锁
//func (r *registry) get(val any) (*model, error) {
// typ := reflect.TypeOf(val)
// r.lock.RLock()
// m, ok := r.models[typ]
// r.lock.RUnlock()
// if ok {
//    return m, nil
// }
//
// r.lock.Lock()
// defer r.lock.Unlock()
// m, ok = r.models[typ]
// if ok {
//    return m, nil
// }
//
// m, err := r.parseModel(val)
// if err != nil {
//    return nil, err
// }
// r.models[typ] = m
// return m, nil
//}

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
