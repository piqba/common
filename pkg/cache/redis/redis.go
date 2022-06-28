// redis
//# Required
//##########
//
//# Set a memory usage limit to the specified amount of bytes.
//# When the memory limit is reached Redis will try to remove keys
//# according to the eviction policy selected (see maxmemory-policy).
//maxmemory 100mb
//
//# Optional
//##########
//
//# Evict any key using approximated LFU when max-memory is reached.
//maxmemory-policy allkeys-lfu
//
//# Enable active memory defragmentation.
//activedefrag yes
//
//# Don't save data on the disk because we can afford to lose cached data.
//save ""

package redis

import (
	"context"
	"fmt"
	ch "github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/piqba/common/pkg/cache"
	"time"
)

// RedisCache ...
type RedisCache struct {
	client *redis.Client
	cache  *ch.Cache
}

// New ...
func New(options cache.CacheOptions) *RedisCache {
	mycache := ch.New(&ch.Options{
		Redis:      options.RdbClient,
		LocalCache: ch.NewTinyLFU(1_000, time.Minute), // pass to cfg
	})
	return &RedisCache{
		client: options.RdbClient,
		cache:  mycache,
	}
}

func (r RedisCache) Set(ctx context.Context, key, payload []byte, duration time.Duration) error {
	if err := r.cache.Set(&ch.Item{
		Ctx:   ctx,
		Key:   string(key),
		Value: payload,
		TTL:   duration,
	}); err != nil {
		return err
	}
	return nil
}

func (r RedisCache) Get(ctx context.Context, key []byte) ([]byte, error) {
	result, err := r.client.Get(ctx, string(key)).Result()
	if err != nil {
		return nil, err
	}
	return []byte(result), nil
}

func (r RedisCache) Invalidate(ctx context.Context, key string) (result int64, err error) {
	iter := r.client.Scan(ctx, 0, fmt.Sprintf("%s:*", key), 0).Iterator()
	for iter.Next(ctx) {
		result, err = r.client.Del(ctx, iter.Val()).Result()
		if err != nil {
			return result, err
		}
	}
	if err := iter.Err(); err != nil {
		return result, err
	}

	return
}
