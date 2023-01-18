package main

import (
	"fmt"
	"time"
)

// 1. go build trace2.go
// 2. GODEBUG=schedtrace=1000 ./trace2    schedtrace=1000是指1000ms
func main() {
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		fmt.Println("Hello GMP")
	}
}
