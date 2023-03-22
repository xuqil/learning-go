package channel

import "context"

type Task func()

type TaskPool struct {
	tasks chan Task
	close chan struct{}
}

// NewTaskPool 新建一个 *TaskPool
// numG 为 goroutine 活动数，capacity 为任务数
func NewTaskPool(numG, capacity int) *TaskPool {
	tp := &TaskPool{
		tasks: make(chan Task, capacity),
		close: make(chan struct{}),
	}

	for i := 0; i < numG; i++ {
		go func() {
			for {
				select {
				case <-tp.close:
					return
				case task := <-tp.tasks:
					task()
				}
			}
		}()
	}

	return tp
}

// Submit 往 TaskPool 提交一个任务
func (tp *TaskPool) Submit(ctx context.Context, t Task) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case tp.tasks <- t:
	}
	return nil
}

// Close 关闭 TaskPool
func (tp *TaskPool) Close() error {
	close(tp.close)
	return nil
}
