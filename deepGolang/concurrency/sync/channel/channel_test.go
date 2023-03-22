package channel

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestChannel(t *testing.T) {
	ch := make(chan string)

	go func() {
		res := <-ch
		fmt.Println("收到消息:", res)
		ch <- "world"
	}()

	// 向 channel 发送信息
	ch <- "hello"

	fmt.Println(<-ch)
}

func TestChannelControl(t *testing.T) {
	ch := make(chan struct{})

	go func() {
		<-ch // 阻塞在这
		fmt.Println("任务启动")
		fmt.Println("任务结束")

		ch <- struct{}{} // 告知主 goroutine 任务结束
	}()

	time.Sleep(time.Second)
	ch <- struct{}{}

	<-ch
}

func TestChannelClose_Recv(t *testing.T) {
	//var wg sync.WaitGroup
	//chInt := make(chan int)
	//wg.Add(1)
	//
	//go func() {
	//	wg.Done()
	//	res := <-chInt
	//	fmt.Println("收到消息:", res) // 收到消息: 0
	//}()
	//close(chInt)
	//wg.Wait()

	chInt := make(chan int)
	close(chInt)
	fmt.Printf("接收的值: %d\n", <-chInt) // 接收的值: 0

	chStr := make(chan string)
	close(chStr)
	fmt.Printf("接收的值: %q\n", <-chStr) // 接收的值: ""
}

func TestChannelClose_Send(t *testing.T) {
	//var wg sync.WaitGroup
	//ch := make(chan int)
	//wg.Add(1)
	//
	//go func() {
	//	defer wg.Done()
	//	res := <-ch
	//	fmt.Println("收到消息:", res) // 收到消息: 0
	//	ch <- 2                   // panic: send on closed channel
	//}()
	//close(ch)
	//
	//wg.Wait()

	ch := make(chan struct{})
	close(ch)
	ch <- struct{}{} // panic: send on closed channel
}

func TestChannelClose_MultiClose(t *testing.T) {
	ch := make(chan int, 10)

	close(ch)
	close(ch) // panic: close of closed channel
}

func TestChannelClose_IsClose(t *testing.T) {
	var wg sync.WaitGroup
	ch := make(chan int)
	wg.Add(1)

	go func() {
		defer wg.Done()
		if res, ok := <-ch; !ok {
			fmt.Println("ch 已经被关闭了") // ch 已经被关闭了
		} else {
			fmt.Println("收到消息:", res)
		}
	}()
	close(ch)
	wg.Wait()
}

func TestChannelClose_Range(t *testing.T) {
	var wg sync.WaitGroup
	ch := make(chan int)
	wg.Add(1)

	go func() {
		defer wg.Done()
		for v := range ch {
			fmt.Println(v)
		}
		// 或者
		//for range ch {
		//}
		fmt.Println("ch 已经被关闭了") // ch 已经被关闭了
	}()
	close(ch)
	wg.Wait()
}

func TestChannel_OnceToOnce(t *testing.T) {
	fmt.Println("start a worker...")
	c := do(func() {
		fmt.Println("worker is working...")
		time.Sleep(time.Second)
	})
	<-c
	fmt.Println("worker work done!")
}

type signal struct{}

func do(f func()) <-chan signal {
	c := make(chan signal)
	go func() {
		fmt.Println("worker start to work...")
		f()
		c <- signal{}
	}()
	return c
}

func TestChannel_OnceToN(t *testing.T) {
	var workers []worker
	for i := 0; i < 10; i++ {
		workers = append(workers, func(i int) {
			fmt.Printf("worker-%d is working...\n", i)
			time.Sleep(time.Second)
		})
	}
	startSingle := make(chan signal)
	done := doMulti(workers, startSingle)

	fmt.Println("start a group of workers...")
	close(startSingle) // 利用 channel close 后的“广播”机制，通知所有 worker 开始工作
	<-done             // 任务完成信号
	fmt.Println("the group of workers work done!")
}

type worker func(i int)

// doMulti 同时执行多个 worker，startSingle 为任务启动信号
func doMulti(workers []worker, startSingle <-chan signal) <-chan signal {
	var wg sync.WaitGroup
	c := make(chan signal)

	for i, w := range workers {
		wg.Add(1)
		go func(f worker, i int) {
			defer wg.Done()
			<-startSingle
			f(i) // 执行任务
		}(w, i)
	}

	go func() {
		wg.Wait()
		c <- signal{}
	}()

	return c
}

func TestCloseNilChannel(t *testing.T) {
	var ch chan int
	close(ch)
}
