package ch12

import (
	"fmt"
	"io"
	"leanring-go/base/ch12/types"
	"os"
	"reflect"
	"testing"
)

func TestReflect_Kind(t *testing.T) {
	var age int8 = 18
	typ := reflect.TypeOf(age)
	switch typ.Kind() {
	case reflect.Int8:
		fmt.Println("age的类型:", reflect.Int8)
	default:
		fmt.Println("age的类型:", typ.Kind().String())
	}
}

func TestReflect_TypeOf(t *testing.T) {
	typ := reflect.TypeOf(3)  // a reflect.Type
	fmt.Println(typ.String()) // "int"
	fmt.Println(typ)          // "int"

	typ2 := reflect.TypeOf(int8(3)) // a reflect.Type
	fmt.Println(typ2.String())      // "int8"
	fmt.Println(typ2)               // "int8"

	var w io.Writer = os.Stdout
	fmt.Println(reflect.TypeOf(w)) // "*os.File"
}

func TestReflect_ValueOf(t *testing.T) {
	//val := reflect.ValueOf(3) // a reflect.Value
	//fmt.Println(val)          // "3"
	//fmt.Printf("%val\n", val)   // "3"
	//fmt.Println(val.String()) // "<int Value>"
	//
	//val2 := reflect.ValueOf(int8(3)) // a reflect.Value
	//fmt.Println(val2)                // "3"
	//fmt.Printf("%val\n", val2)         // "3"
	//fmt.Println(val2.String())       // "<int8 Value>"

	//val := reflect.ValueOf(3) // a reflect.Value
	//x := val.Interface()      // an interface{}
	//i := x.(int)              // an int
	//fmt.Printf("%d\n", i)     // "3"
	//
	//val2 := reflect.ValueOf(int8(3)) // a reflect.Value
	//x2 := val2.Interface()           // an interface{}
	//i2 := x2.(int8)                  // an int8
	//fmt.Printf("%d\n", i2)           // "3"

	val := reflect.ValueOf(3) // a reflect.Value
	typ := val.Type()         // a reflect.Type
	fmt.Println(typ.String()) // "int"

	val2 := reflect.ValueOf(int8(3)) // a reflect.Value
	typ2 := val2.Type()              // a reflect.Type
	fmt.Println(typ2.String())       // "int8"

}

func TestReflect_String(t *testing.T) {
	var name = "Jerry"
	nTyp := reflect.TypeOf(name)
	nVal := reflect.ValueOf(name)

	// string->type: string value: Jerry
	fmt.Printf("%s->type: %s value: %s\n", nTyp.String(), nTyp.Kind().String(), nVal.Interface())
}

func TestReflect_StringPtr(t *testing.T) {
	var name = "Jerry"

	nTyPtr := reflect.TypeOf(&name)
	nValPtr := reflect.ValueOf(&name)

	// *string->type: ptr value: 0xc0000745e0
	fmt.Printf("%s->type: %s value: %v\n", nTyPtr.String(), nTyPtr.Kind().String(), nValPtr.Interface())
	// *string->type: ptr value: Jerry
	fmt.Printf("%s->type: %s value: %v\n", nTyPtr.String(), nTyPtr.Kind().String(), nValPtr.Elem().Interface())
}
func TestReflect_StringExporter(t *testing.T) {
	var ms types.Message = "hello, world"
	nTyp := reflect.TypeOf(ms)
	nVal := reflect.ValueOf(ms)

	// string->type: string value: Jerry
	fmt.Printf("%s->type: %s value: %s\n", nTyp.String(), nTyp.Kind().String(), nVal.Interface())
}

func TestReflect_StructField(t *testing.T) {

	type ReflectUser struct {
		Name string
		Age  int
		// 如果同属一个包，phone 可以被测试访问到，如果是不同包，就访问不到了
		phone string
	}

	ru := ReflectUser{
		Name:  "Tom",
		Age:   18,
		phone: "666",
	}

	IterateFields(&ru)
}

func TestReflect_StructPtrField(t *testing.T) {

	type ReflectUser struct {
		Name string
		Age  int
		// 如果同属一个包，phone 可以被测试访问到，如果是不同包，就访问不到了
		phone string
	}

	ru := ReflectUser{
		Name:  "Tom",
		Age:   18,
		phone: "666",
	}

	//IterateFields(&ru)
	typ := reflect.TypeOf(ru)
	ft, found := typ.FieldByName("Name")
	if found {
		fmt.Println(ft.Name)
	}
}

func TestReflect_SetValue(t *testing.T) {
	var name = "Tom"
	fmt.Println("name before:", name)
	val := reflect.ValueOf(&name)
	val = val.Elem()
	fmt.Println(val.Type().Kind()) // "ptr"
	if !val.CanSet() {
		fmt.Println("Name不可被修改")
	} else {
		val.Set(reflect.ValueOf("Jerry"))
	}
	fmt.Println("name after:", name)
}

type ReflectUser struct {
	Name  string
	Age   int
	phone string
}

func (r ReflectUser) String() string {
	return fmt.Sprintf("Name: %s Age: %d phone: %s",
		r.Name, r.Age, r.phone)
}
func TestReflect_StructSetValue(t *testing.T) {
	ru := ReflectUser{
		Name:  "Tom",
		Age:   18,
		phone: "666",
	}
	fmt.Println("before:", ru)
	SetField(&ru, "Name", "Jerry")
	SetField(&ru, "Age", 20)
	SetField(&ru, "phone", "888") // unexported 字段不可修改
	fmt.Println("after:", ru)
}

func TestReflect_Method(t *testing.T) {
	//u := types.NewUser("Tom", 18, "666")
	//f := types.User.GetAge
	//fmt.Printf("%T\n", u.GetAge)             // func() int
	//fmt.Printf("%T\n", f)                    // func(types.User) int
	//fmt.Printf("%T\n", (*types.User).GetAge) //func(*types.User) int
	//
	//// 下面两条语句是等价的
	//fmt.Println(u.GetAge()) // 18
	//fmt.Println(f(u))       // 18

	user := types.NewUser("Tom", 18, "666")
	//user := types.NewUserPtr("Tom", 18, "666")
	IterateFunc(user)
}

func TestIterateFunc(t *testing.T) {
	//user := types.NewUser("Tom", 18, "666")
	user := types.NewUserPtr("Tom", 18, "666")
	IterateFunc(user)
}
