package cache

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/go-redis/redis/v8"
)

const (
	Memory = "memory"
	Redis  = "redis"
)

type CacheOptions struct {
	RdbClient      *redis.Client
	BadgerClient   *badger.DB
	BadgerInMemory bool
	// #example "./badger.data"
	BadgerDiskPath string

	// Cacheable Criteria
	CacheableCriteria interface{}
}
