package waitGroup

import (
	"fmt"
	"sync"
	"testing"
)

func TestWaiGroup(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		// 每增加一个 goroutine 都要调用 Add 加 1
		wg.Add(1)
		go func(i int) {
			// goroutine 执行完一定要调用 Done，即 Add(-1)
			defer wg.Done()
			fmt.Println("task", i, "done")
		}(i)
	}
	// 主 goroutine 等待所有子 goroutine 完成
	wg.Wait()
	fmt.Println("all task done")
}
