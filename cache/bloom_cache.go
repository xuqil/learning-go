package cache

import (
	"context"
	"fmt"
)

type BloomFilterCache struct {
	ReadThroughCache
}

func NewBloomFilterCache(cache Cache, bf BloomFilter,
	loadFunc func(ctx context.Context, key string) (any, error)) *BloomFilterCache {
	return &BloomFilterCache{
		ReadThroughCache{
			Cache: cache,
			LoadFunc: func(ctx context.Context, key string) (any, error) {
				if !bf.HasKey(ctx, key) {
					return nil, errKeyNotFound
				}
				return loadFunc(ctx, key)
			},
		},
	}
}

type BloomFilterCacheV1 struct {
	ReadThroughCache
	Bf BloomFilter
}

func (b *BloomFilterCacheV1) Get(ctx context.Context, key string) (any, error) {
	val, err := b.Cache.Get(ctx, key)
	if err == errKeyNotFound && b.Bf.HasKey(ctx, key) {
		val, err = b.LoadFunc(ctx, key)
		if err == nil {
			er := b.Cache.Set(ctx, key, val, b.Expiration)
			if err != nil {
				return val, fmt.Errorf("%w, 原因：%s", ErrFailedToRefreshCache, er.Error())
			}
		}
	}
	return val, err
}

type BloomFilter interface {
	HasKey(ctx context.Context, key string) bool
}
