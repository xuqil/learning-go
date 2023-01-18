package ch12

import (
	"fmt"
	"reflect"
)

func IterateFields(entity any) {
	typ := reflect.TypeOf(entity)
	val := reflect.ValueOf(entity)
	if val.IsZero() {
		fmt.Println("不支持零值结构体")
		return
	}
	for typ.Kind() == reflect.Ptr {
		// 拿到指针指向的对象
		typ = typ.Elem()
		val = val.Elem()
	}
	if typ.Kind() != reflect.Struct {
		fmt.Println("不是结构体类型")
		return
	}
	for i := 0; i < typ.NumField(); i++ {
		ft := typ.Field(i)
		fv := val.Field(i)
		if ft.IsExported() {
			fmt.Printf("%s.%s 的类型: %s  值: %v\n",
				typ.Name(), ft.Name, ft.Type.String(), fv.Interface())
		} else {
			fmt.Printf("%s.%s 的类型: %s  值: %q\n",
				typ.Name(), ft.Name, ft.Type.String(), reflect.Zero(ft.Type).Interface())
		}
	}
}

func SetField(entity any, field string, newValue any) {
	val := reflect.ValueOf(entity)
	for val.Type().Kind() == reflect.Pointer {
		val = val.Elem()
	}
	fv := val.FieldByName(field)
	if !fv.CanSet() {
		fmt.Println(field, "不可修改字段")
		return
	}
	fv.Set(reflect.ValueOf(newValue))
}
