package cache

import (
	"context"
	"fmt"
	"golang.org/x/sync/singleflight"
	"time"
)

type SingleflightCacheV1 struct {
	ReadThroughCache
}

// NewSingleflightCacheV1 新建一个 singlefilght ReadThoughCache
// 读 DB 时只有一个 goroutine，但回写 Cache 时会有多个 goroutine
func NewSingleflightCacheV1(cache Cache,
	loadFunc func(ctx context.Context, key string) (any, error),
	expiration time.Duration) *SingleflightCacheV1 {
	g := &singleflight.Group{}
	return &SingleflightCacheV1{
		ReadThroughCache{
			Cache: cache,
			LoadFunc: func(ctx context.Context, key string) (any, error) {
				val, err, _ := g.Do(key, func() (interface{}, error) {
					return loadFunc(ctx, key)
				})
				return val, err
			},
			Expiration: expiration,
		},
	}
}

// SingleflightCacheV2 singlefilght ReadThoughCache
// 读 DB 和回写 Cache 时只有一个 goroutine
type SingleflightCacheV2 struct {
	ReadThroughCache
	g *singleflight.Group
}

func (r *SingleflightCacheV2) Get(ctx context.Context, key string) (any, error) {
	val, err := r.Cache.Get(ctx, key)
	if err == errKeyNotFound {
		val, err, _ = r.g.Do(key, func() (interface{}, error) {
			v, er := r.LoadFunc(ctx, key)
			if er == nil {
				//_ = r.Cache.Set(ctx, key, val, r.Expiration)
				er = r.Cache.Set(ctx, key, val, r.Expiration)
				if er != nil {
					return v, fmt.Errorf("%w, 原因：%s", ErrFailedToRefreshCache, er.Error())
				}
			}
			return v, er
		})
	}
	return val, err
}
