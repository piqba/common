package factory

import (
	"context"
	"errors"
	"github.com/piqba/common/pkg/cache"
	"github.com/piqba/common/pkg/cache/badgerdb"
	"github.com/piqba/common/pkg/cache/redis"
	"time"
)

type ICache interface {
	Set(ctx context.Context, key, payload []byte, duration time.Duration) error
	Get(ctx context.Context, key []byte) ([]byte, error)
	Invalidate(ctx context.Context, key string) (int64, error)
}

// GetCacheFactory it`s a constructor factory method
func GetCacheFactory(driverType string, options cache.CacheOptions) (ICache, error) {
	switch driverType {
	case cache.Memory:
		return badgerdb.New(options), nil
	case cache.Redis:
		return redis.New(options), nil
	}
	return nil, errors.New("error: Bad repo type ")
}
