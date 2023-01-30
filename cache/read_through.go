package cache

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

var (
	ErrFailedToRefreshCache = errors.New("刷新缓存失败")
)

// ReadThroughCache 你一定要赋值 LoadFunc 和 Expiration
// Expiration 是你的过期时间
type ReadThroughCache struct {
	Cache
	LoadFunc   func(ctx context.Context, key string) (any, error)
	Expiration time.Duration
}

// Get 同步 LoadFunc 和刷新缓存
func (r *ReadThroughCache) Get(ctx context.Context, key string) (any, error) {
	val, err := r.Cache.Get(ctx, key)
	if err == errKeyNotFound {
		val, err = r.LoadFunc(ctx, key)
		if err == nil {
			er := r.Cache.Set(ctx, key, val, r.Expiration)
			if err != nil {
				return val, fmt.Errorf("%w, 原因：%s", ErrFailedToRefreshCache, er.Error())
			}
		}
	}
	return val, err
}

// GetV1 异步 LoadFunc 和刷新缓存
func (r *ReadThroughCache) GetV1(ctx context.Context, key string) (any, error) {
	val, err := r.Cache.Get(ctx, key)
	if err == errKeyNotFound {
		go func() {
			val, err = r.LoadFunc(ctx, key)
			if err == nil {
				//_ = r.Cache.Set(ctx, key, val, r.Expiration)
				er := r.Cache.Set(ctx, key, val, r.Expiration)
				if er != nil {
					log.Fatalln(er)
				}
			}
		}()
	}
	return val, err
}

// GetV2 同步 LoadFunc 但异步刷新缓存
func (r *ReadThroughCache) GetV2(ctx context.Context, key string) (any, error) {
	val, err := r.Cache.Get(ctx, key)
	if err == errKeyNotFound {
		val, err = r.LoadFunc(ctx, key)
		if err == nil {
			go func() {
				//_ = r.Cache.Set(ctx, key, val, r.Expiration)
				er := r.Cache.Set(ctx, key, val, r.Expiration)
				if er != nil {
					log.Fatalln(er)
				}

			}()
		}
	}
	return val, err
}
