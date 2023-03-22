package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx := context.Background()
	timeoutCtx, cancel1 := context.WithTimeout(ctx, time.Second*2)
	subCtx, cancel2 := context.WithTimeout(timeoutCtx, time.Second*1)
	go func() {
		<-subCtx.Done() // subCtx 会在一秒钟后过期，先输出 timeout2
		fmt.Println("timeout2")
	}()
	go func() {
		<-timeoutCtx.Done() // timeoutCtx 会在两秒钟后过期，然后输出 timeout1
		fmt.Println("timeout1")
	}()
	time.Sleep(3 * time.Second)
	cancel2()
	cancel1()
}
