package redis

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/piqba/common/pkg/databases"
)

//Used to execute client creation procedure only once.
var redisOnce sync.Once

// NewRedisDbOnce Create an unique redis instance
func NewRedisDbOnce(ctx context.Context, option databases.RdbOptions) (*redis.Client, error) {
	var clientInstance *redis.Client
	var err error
	redisOnce.Do(func() {

		client := redis.NewClient(&redis.Options{
			Addr:         option.Addr,
			Username:     option.Username,
			Password:     option.Password,
			DB:           option.DB,
			DialTimeout:  60 * time.Second,
			ReadTimeout:  60 * time.Second,
			WriteTimeout: 60 * time.Second,
		})

		_, errPing := client.Ping(ctx).Result()
		if errPing != nil {
			err = errPing
		}
		clientInstance = client
	})

	return clientInstance, err
}

// NewRedisDb Create a redis instance
func NewRedisDb(ctx context.Context, option databases.RdbOptions) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         option.Addr,
		Username:     option.Username,
		Password:     option.Password,
		DB:           option.DB,
		DialTimeout:  60 * time.Second,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return client, err
}
