package ch12

import (
	"fmt"
	"reflect"
)

func IterateFunc(entity any) {
	typ := reflect.TypeOf(entity)
	numMethod := typ.NumMethod()
	fmt.Println(typ.String(), "方法的个数:", numMethod)
	for i := 0; i < numMethod; i++ {
		method := typ.Method(i)
		fn := method.Func // 拿到结构体方法的 Value，需要注意的是，方法的第一个参数为 receiver

		numIn := fn.Type().NumIn()                     // 方法参数的个数
		inputTypes := make([]reflect.Type, 0, numIn)   // 每个参数的类型
		inputValues := make([]reflect.Value, 0, numIn) //每个入参的零值

		inputValues = append(inputValues, reflect.ValueOf(entity)) // 第一个参数为 receiver
		inputTypes = append(inputTypes, reflect.TypeOf(entity))    // 第一个入参为 receiver

		paramTypes := fmt.Sprintf("%s", typ.String())
		// 这个遍历是为了得到方法的各个参数类型，以及参数对应的零值
		for j := 1; j < numIn; j++ {
			fnInType := fn.Type().In(j)
			inputTypes = append(inputTypes, fnInType)                 // append 参数的类型
			inputValues = append(inputValues, reflect.Zero(fnInType)) // append 入参的零值

			paramTypes = paramTypes + "," + fnInType.String()
		}

		returnTypes := ""

		numOut := fn.Type().NumOut() // 返回值的个数
		outputTypes := make([]reflect.Type, 0, numOut)
		// 拿到每个返回值的类型
		for j := 0; j < numOut; j++ {
			fnOutType := fn.Type().Out(j)
			outputTypes = append(outputTypes, fnOutType)

			if j > 0 {
				returnTypes += ","
			}
			returnTypes += fnOutType.String()
		}

		resValues := fn.Call(inputValues) // 调用结构体里的方法，Call 的参数为方法入参的切片，返回值存储在切片
		result := make([]any, 0, len(resValues))
		for _, v := range resValues {
			result = append(result, v.Interface())
		}

		fSign := fmt.Sprintf("func(%s) %s", paramTypes, returnTypes)
		fmt.Println("方法签名:", fSign)
		fmt.Printf("调用方法: %s 入参: %v 返回结果: %v\n", method.Name, inputValues, result)
	}
}
