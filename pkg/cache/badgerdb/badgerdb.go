package badgerdb

import (
	"context"
	badger "github.com/dgraph-io/badger/v3"
	"github.com/piqba/common/pkg/cache"
	"time"
)

type BadgerCache struct {
	clientBadger *badger.DB
}

// New ...
func New(options cache.CacheOptions) *BadgerCache {
	return &BadgerCache{
		clientBadger: options.BadgerClient,
	}
}

func (b *BadgerCache) Set(_ context.Context, key []byte, payload []byte, duration time.Duration) error {
	if err := b.clientBadger.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(key, payload).WithTTL(duration)
		err := txn.SetEntry(e)
		return err
	}); err != nil {
		return err
	}
	return nil
}

func (b *BadgerCache) Get(_ context.Context, key []byte) ([]byte, error) {
	var resultPayload []byte
	if err := b.clientBadger.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		if err := item.Value(func(value []byte) error {
			resultPayload = value
			return nil
		}); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return resultPayload, nil
}

func (b *BadgerCache) Invalidate(ctx context.Context, key string) (int64, error) {
	if err := b.clientBadger.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		err = txn.Delete([]byte(key))
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return 0, err
	}
	return 0, nil
}
