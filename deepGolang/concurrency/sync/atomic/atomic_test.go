package atomic

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

func TestNoAtomic(t *testing.T) {
	res := 0
	target := 1000
	var wg sync.WaitGroup
	for i := 1; i <= target; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			res += i
		}(i)
	}
	wg.Wait()
	t.Log("预期结果:", sum(target))
	t.Log("实际结果:", res)
}

func TestNoAtomic_OneGoroutine(t *testing.T) {
	runtime.GOMAXPROCS(1)

	res := 0
	target := 1000
	var wg sync.WaitGroup
	for i := 1; i <= target; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			res += i
		}(i)
	}
	wg.Wait()
	t.Log("预期结果:", sum(target))
	t.Log("实际结果:", res)
}

func TestSerial(t *testing.T) {
	res := 0
	target := 1000
	for i := 1; i <= target; i++ {
		res = res + i
	}
	t.Log("预期结果:", sum(target)) //预期结果: 500500
	t.Log("实际结果:", res)         //实际结果: 500500
}

func TestAtomic(t *testing.T) {
	var res int32
	target := 1000
	var wg sync.WaitGroup
	for i := 1; i <= target; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			atomic.AddInt32(&res, int32(i))
		}(i)
	}
	wg.Wait()
	t.Log("预期结果:", sum(target)) //预期结果: 500500
	t.Log("实际结果:", res)         //实际结果: 500500
}

func TestMutex(t *testing.T) {
	res := 0
	target := 1000
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i := 1; i <= target; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			mu.Lock()
			res += i
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	t.Log("预期结果:", sum(target)) //预期结果: 500500
	t.Log("实际结果:", res)         //实际结果: 500500
}

func sum(t int) int {
	return (1 + t) * t / 2
}

func TestAddX(t *testing.T) {
	var a uint32 = 10
	atomic.AddUint32(&a, ^uint32(2-1))
	t.Log(a)
}

func TestCAS(t *testing.T) {
	var a int64 = 9
	t.Log("a before:", a)
	t.Log("a交换是否成功", atomic.CompareAndSwapInt64(&a, 10, 11))
	t.Log("a after:", a)

	var b int64 = 10
	t.Log("b before:", b)
	t.Log("b交换是否成功", atomic.CompareAndSwapInt64(&b, 10, 11))
	t.Log("b after:", b)
}

func TestSwap(t *testing.T) {
	var a int32 = 9
	atomic.SwapInt32(&a, 10)
	t.Log(a)
}

func TestStoreAndLoad(t *testing.T) {
	var a int32

	atomic.StoreInt32(&a, 10)
	res := atomic.LoadInt32(&a)
	t.Log(res) // 10
}

func TestValue(t *testing.T) {
	var a atomic.Value
	a.Store(int32(10))
	res := a.Load().(int32)
	t.Log(res) // 10

	t.Log(a.CompareAndSwap(int32(10), int32(11))) // true
	res = a.Load().(int32)
	t.Log(res) // 11
}
