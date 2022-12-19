package main

import (
	"fmt"
	"os"
	"runtime/trace"
)

// trace 的编程过程:
// 1. 创建文件
// 2. 启动
// 3. 停止

// 打开和解析 trace 文件：go tool trace trace.out
func main() {
	// 1. 创建一个 trace 文件
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	// 2. 启动 trace
	err = trace.Start(f)
	if err != nil {
		panic(err)
	}

	// 正常调试的业务
	fmt.Println("Hello GMP")

	// 3. 停止 trace
	trace.Stop()
}
