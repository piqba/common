package badgerdb

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/piqba/common/pkg/cache"
)

func NewBadgerDB(options cache.CacheOptions) (*badger.DB, error) {
	var opt badger.Options
	if options.BadgerInMemory {
		opt = badger.DefaultOptions("").
			WithInMemory(true).
			WithBypassLockGuard(true)
	} else {
		opt = badger.DefaultOptions(options.BadgerDiskPath).
			WithBypassLockGuard(true)
	}

	clientBadger, err := badger.Open(opt)
	if err != nil {
		return nil, err
	}

	return clientBadger, nil
}
