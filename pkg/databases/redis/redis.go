package redis

import (
	"context"
	"github.com/piqba/common/pkg/databases"
	"runtime"
	"sync"
	"time"
)

//Used to execute client creation procedure only once.
var redisOnce sync.Once

// NewRedisDbOnce Create an unique redis instance
func NewRedisDbOnce(ctx context.Context, option databases.RdbOptions) (*redis.Client, error) {
	var clientInstance *redis.Client
	redisOnce.Do(func() {

		client := redis.NewClient(&redis.Options{
			Addr: option.Addr,
			OnConnect: func(ctx context.Context, cn *redis.Conn) error {
				return cn.Ping(ctx).Err()
			},
			Username:        option.Username,
			Password:        option.Password,
			DB:              option.DB,
			MaxRetries:      10,
			MinRetryBackoff: 5,
			MaxRetryBackoff: 5,
			DialTimeout:     60 * time.Second,
			ReadTimeout:     60 * time.Second,
			WriteTimeout:    60 * time.Second,
			PoolSize:        runtime.GOMAXPROCS(runtime.NumCPU()),
			MinIdleConns:    50,
		})

		clientInstance = client
	})

	return clientInstance, nil
}

// NewRedisDb Create a redis instance
func NewRedisDb(ctx context.Context, option databases.RdbOptions) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: option.Addr,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			return cn.Ping(ctx).Err()
		},
		Username:        option.Username,
		Password:        option.Password,
		DB:              option.DB,
		MaxRetries:      10,
		MinRetryBackoff: 5,
		MaxRetryBackoff: 5,
		DialTimeout:     60 * time.Second,
		ReadTimeout:     60 * time.Second,
		WriteTimeout:    60 * time.Second,
		PoolSize:        runtime.GOMAXPROCS(runtime.NumCPU()),
		MinIdleConns:    50,
	})

	return client, nil
}
