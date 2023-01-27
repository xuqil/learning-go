package context

import (
	"context"
	"testing"
	"time"
)

type mykey struct{}
type mykeyv2 int

// Background 没有 timeout 和 cancel 的 context
func TestContext(t *testing.T) {
	// 一般是链路起点，或者调用的起点，用 Background()
	ctx := context.Background()
	// 不确定 context 该用啥的时候，用 TODO()
	//ctx := context.TODO()

	ctx = context.WithValue(ctx, mykey{}, "my-value")
	val := ctx.Value(mykey{}).(string)
	t.Log(val)

	newVal := ctx.Value("不存在的 key")
	val, ok := newVal.(string)
	if !ok {
		t.Log("类型不对")
		return
	}
	t.Log(val)

	//t.Log(ctx.Value("不存在的 key").(string)) // 连调会 panic
}

// WithCancel 带 cancel 的 context
func TestContext_WithCancel(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	// 用完 ctx 再调用 cancel()
	//defer cancel()

	go func() {
		time.Sleep(time.Second)
		cancel()
	}()

	// 用 ctx
	<-ctx.Done() // cancel() 执行后才不会阻塞
	t.Log("hello, cancel: ", ctx.Err())
}

// WithDeadline 带 timeout 和 cancel 的context
func TestContext_WithDeadline(t *testing.T) {
	ctx := context.Background()
	// 3 秒后自动 deadline
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	deadline, _ := ctx.Deadline()
	t.Log("deadline: ", deadline)
	defer cancel()
	<-ctx.Done()
	t.Log("hello, deadline: ", ctx.Err())
}

// WithTimeout 是 WithDeadline 的封装
func TestContext_WithTimeout(t *testing.T) {
	ctx := context.Background()
	// 3 秒后自动 deadline
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	deadline, _ := ctx.Deadline()
	t.Log("deadline: ", deadline)
	defer cancel()
	<-ctx.Done()
	t.Log("hello, timeout: ", ctx.Err())
}

func TestContext_Parent(t *testing.T) {
	ctx := context.Background()
	parent := context.WithValue(ctx, "my-key", "my value")
	child := context.WithValue(parent, "my-key", "my new value")

	t.Log("parent my-key: ", parent.Value("my-key")) // parent my-key:  my value
	t.Log("child my-key: ", child.Value("my-key"))   // child my-key:  my new value

	child2, cancel := context.WithTimeout(parent, time.Second)
	defer cancel()
	t.Log("child2 my-key: ", child2.Value("my-key")) // child2 my-key:  my value

	// 父 context 不能访问子 context 的 value
	child3 := context.WithValue(parent, "new-key", "child3 value")
	t.Log("parent new-key: ", parent.Value("new-key")) // parent new-key:  <nil>
	t.Log("child3 new-key: ", child3.Value("new-key")) // child3 new-key:  child3 value

	// 逼不得已使用（父 context 传一个可变【不可哈希】的“筐”给子 context，子 context 往“筐”里装东西，父 context 是可以访问到“筐”里的东西）
	parent1 := context.WithValue(parent, "map", map[string]string{})
	child4, cancel := context.WithTimeout(parent1, time.Second)
	defer cancel()
	m := child4.Value("map").(map[string]string)
	m["key1"] = "value1"
	nm := parent1.Value("map").(map[string]string)
	t.Log("parent1 key1: ", nm["key1"]) // parent1 key1:  value1
}
