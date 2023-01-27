package cache

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

var (
	errOverCapacity = errors.New("cache：超过容量限制")
)

type MaxCntCache struct {
	*BuildInMapCache
	cnt    int32
	maxCnt int32
}

func NewMaxCntCache(c *BuildInMapCache, maxCnt int32) *MaxCntCache {
	res := &MaxCntCache{
		BuildInMapCache: c,
		maxCnt:          maxCnt,
	}

	origin := c.onEvicted
	c.onEvicted = func(key string, val any) {
		atomic.AddInt32(&res.cnt, -1)
		if origin != nil {
			origin(key, val)
		}
	}

	return res
}

func (c *MaxCntCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, ok := c.data[key]
	if !ok {
		if c.cnt+1 > c.maxCnt {
			return errOverCapacity
		}
		atomic.AddInt32(&c.cnt, 1)
	}

	return c.set(key, val, expiration)
}
