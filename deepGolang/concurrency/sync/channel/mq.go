package channel

import (
	"errors"
	"sync"
)

type Msg struct {
	Content string
}

type Broker struct {
	mutex sync.RWMutex
	chans []chan Msg
}

// Send 向每个 channel 发送消息
func (b *Broker) Send(m Msg) error {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	for _, ch := range b.chans {
		select {
		case ch <- m:
		default:
			return errors.New("消息队列已满")
		}
	}
	return nil
}

// Subscribe 订阅一个 capacity 大小的队列
func (b *Broker) Subscribe(capacity int) (<-chan Msg, error) {
	ch := make(chan Msg, capacity)
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.chans = append(b.chans, ch)
	return ch, nil
}

// Close 关闭 Broker
func (b *Broker) Close() error {
	b.mutex.Lock()
	chans := b.chans
	b.chans = nil
	b.mutex.Unlock()

	for _, ch := range chans {
		close(ch)
	}
	return nil
}
