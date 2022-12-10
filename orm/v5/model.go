package orm

import (
	"leanring-go/orm/internal/errs"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

const (
	tagKeyColumn = "column"
)

type Registry interface {
	Get(val any) (*Model, error)
	Register(val any, opts ...ModelOpt) (*Model, error)
}

type Model struct {
	tableName string
	fields    map[string]*Field
}

type ModelOpt func(m *Model) error

type Field struct {
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

// Get 会存在覆盖问题
func (r *registry) Get(val any) (*Model, error) {
	typ := reflect.TypeOf(val)
	m, ok := r.models.Load(typ)
	if ok {
		return m.(*Model), nil
	}
	m, err := r.Register(val)
	if err != nil {
		return nil, err
	}
	return m.(*Model), nil
}

// Get 使用 double read 读写锁
//func (r *registry) Get(val any) (*model, error) {
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
// m, err := r.Register(val)
// if err != nil {
//    return nil, err
// }
// r.models[typ] = m
// return m, nil
//}

func (r *registry) Register(entity any, opts ...ModelOpt) (*Model, error) {
	typ := reflect.TypeOf(entity)
	// 限制只能用一级指针
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	elemType := typ.Elem()
	numField := elemType.NumField()
	fieldMap := make(map[string]*Field, numField)
	for i := 0; i < numField; i++ {
		fd := elemType.Field(i)
		pair, err := r.parseTag(fd.Tag)
		if err != nil {
			return nil, err
		}
		colName := pair[tagKeyColumn]
		if colName == "" {
			// 用户没有设置
			colName = underscoreName(fd.Name)
		}
		fieldMap[fd.Name] = &Field{
			colName: colName,
		}
	}
	var tableName string
	if tbl, ok := entity.(TableName); ok {
		tableName = tbl.TableName()
	}

	if tableName == "" {
		tableName = underscoreName(elemType.Name())
	}

	res := &Model{
		tableName: tableName,
		fields:    fieldMap,
	}

	for _, opt := range opts {
		err := opt(res)
		if err != nil {
			return nil, err
		}
	}
	r.models.Store(typ, res)
	return res, nil
}

func ModelWithTableName(tableName string) ModelOpt {
	return func(m *Model) error {
		m.tableName = tableName
		//if tableName == "" {
		//	return errs
		//}
		return nil
	}
}

func ModelWithColumnName(field string, colName string) ModelOpt {
	return func(m *Model) error {
		fd, ok := m.fields[field]
		if !ok {
			return errs.NewErrUnknownField(field)
		}
		fd.colName = colName
		return nil
	}
}

func (r *registry) parseTag(tag reflect.StructTag) (map[string]string, error) {
	ormTag, ok := tag.Lookup("orm")
	if !ok {
		return map[string]string{}, nil
	}
	pairs := strings.Split(ormTag, ",")
	res := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		segs := strings.Split(pair, "=")
		if len(segs) != 2 {
			return nil, errs.NewErrInvalidTagContent(pair)
		}
		key := segs[0]
		val := segs[1]
		res[key] = val
	}
	return res, nil
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