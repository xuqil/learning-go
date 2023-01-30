package cache

import (
	"context"
	"log"
	"time"
)

type WriteThroughCache struct {
	Cache
	StoreFunc func(ctx context.Context, key string, val any) error
}

// Set 先写 DB 再写 Cache
func (w *WriteThroughCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	err := w.StoreFunc(ctx, key, val)
	if err != nil {
		return err
	}
	return w.Cache.Set(ctx, key, val, expiration)
}

// SetV1 先写 Cache 再写 DB
func (w *WriteThroughCache) SetV1(ctx context.Context, key string, val any, expiration time.Duration) error {
	err := w.Cache.Set(ctx, key, val, expiration)
	if err != nil {
		return err
	}
	return w.StoreFunc(ctx, key, val)
}

// SetV2 先写 DB 再异步写 Cache
func (w *WriteThroughCache) SetV2(ctx context.Context, key string, val any, expiration time.Duration) error {
	err := w.StoreFunc(ctx, key, val)
	if err == nil {
		go func() {
			er := w.Cache.Set(ctx, key, val, expiration)
			if er != nil {
				log.Fatalln(er)
			}
		}()
	}
	return err
}
