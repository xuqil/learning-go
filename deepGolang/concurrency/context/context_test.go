package main

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestTimeoutExample(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	bsChan := make(chan struct{}) // 业务的 channel
	go func() {
		bs()
		bsChan <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("timeout")
	case <-bsChan:
		fmt.Println("business end")
	}
}

func TestTimeoutTimeAfterFunc(t *testing.T) {
	bsChan := make(chan struct{})
	go func() {
		bs()
		bsChan <- struct{}{}
	}()

	timer := time.AfterFunc(time.Second, func() {
		fmt.Println("timeout")
	})
	<-bsChan
	timer.Stop() // 取消 timer
}
func bs() {
	time.Sleep(time.Second * 2)
}

func fun() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	<-ctx.Done()
}

func TestContext(t *testing.T) {
	ctx := Context{
		ch: make(chan struct{}),
	}
	go func() {
		for {
			select {
			case <-ctx.done():
				fmt.Println("cancel")
				return
			}
		}
	}()

	go func() {
		time.Sleep(time.Second * 1)
		ctx.cancel()
	}()

	time.Sleep(time.Second * 2)
	fmt.Println("exit")
}

type Context struct {
	ch chan struct{}
}

func (c *Context) cancel() {
	c.ch <- struct{}{}
	//close(c.ch)
}

func (c *Context) done() <-chan struct{} {
	return c.ch
}

func TestCancelContext(t *testing.T) {
	ctx1, cancel1 := context.WithCancel(context.Background())
	ctx, cancel := context.WithCancel(ctx1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			}
		}
	}()

	time.Sleep(time.Second * 2)
	cancel1()
	cancel()
}
