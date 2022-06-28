# Cache module

This module, pkg or folder contain two possible datasource `redis` & [`badgerdb`](https://github.com/dgraph-io/badger/)

>Cache types
- Least Recently Used (`LRU`): remove an entry that has been used least recently.
- First In, First Out (`FIFO`): remove an entry that was created first.
- Least Frequently Used (`LFU`): remove an entry that was least frequently used.

it's used the common design pattern `factory`

```go
type ICache interface {
	Set(ctx context.Context, key, payload []byte, duration time.Duration) error
	Get(ctx context.Context, key []byte) ([]byte, error)
}

// GetCacheFactory it`s a constructor factory method
func GetCacheFactory(driverType string, options CacheOptions) (ICache, error) {
	switch driverType {
	case Memory:
		return badgerdb.transformer(options), nil
	case Redis:
		return redis.transformer(options), nil
	}
	return nil, errors.transformer("error: Bad repo type ")
}
```

For options definitions consult `types.go` file

```go

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
}

```

### TODO
- Implementar estrategia de cache (consultar)