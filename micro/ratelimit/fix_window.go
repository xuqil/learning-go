package ratelimit

import (
	"context"
	"google.golang.org/grpc"
	"sync/atomic"
	"time"
)

type FixWindowLimiter struct {
	// 窗口的起始时间
	timestamp int64
	// 窗口大小
	interval int64
	// 在这个窗口内，允许通过的最大请求数量
	rate     int64
	cnt      int64
	onReject rejectStrategy
}

func NewFixWindowLimiter(interval time.Duration, rate int64) *FixWindowLimiter {
	return &FixWindowLimiter{
		interval:  interval.Nanoseconds(),
		timestamp: time.Now().UnixNano(),
		rate:      rate,
		onReject:  defaultRejectStrategy,
	}
}

func (t *FixWindowLimiter) BuildServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// 考虑 t.cnt 重置的问题
		current := time.Now().UnixNano()
		timestamp := atomic.LoadInt64(&t.timestamp)
		cnt := atomic.LoadInt64(&t.cnt)
		if timestamp+t.interval < current {
			// 新窗口，需要重置窗口
			if atomic.CompareAndSwapInt64(&t.timestamp, timestamp, current) {
				atomic.CompareAndSwapInt64(&t.cnt, cnt, 0)
			}
		}
		cnt = atomic.AddInt64(&t.cnt, 1)
		if cnt > t.rate {
			//err = errors.New("触发瓶颈了")
			//return
			return t.onReject(ctx, req, info, handler)
		}
		resp, err = handler(ctx, req)
		return
	}
}
