package pool

import (
	"sync"
	"testing"
)

type Student struct {
	Name string
}

func TestPool(t *testing.T) {
	p := sync.Pool{
		New: func() any {
			//	创建函数， sync.Pool 会回调，用于创建对象
			t.Log("创建一个对象")
			return Student{}
		},
	}
	for i := 0; i < 10; i++ {
		// 从 Pool 中取出对象
		// 1. 如果 Pool 没有，会调用 New() 创建一个然后返回
		// 2. 如果 Pool 有对象，则直接返回
		obj := p.Get().(Student)

		// 用完对象后需要放回去
		p.Put(obj)
	}
}
