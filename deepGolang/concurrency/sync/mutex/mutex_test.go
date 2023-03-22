package mutex

import "sync"

// safeResource 所有对资源的操作都只能通过定义在 safeResource 上的方法来进行
type safeResource struct {
	resource any
	mu       sync.Mutex
}
