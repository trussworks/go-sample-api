package cache

import (
	"time"

	"github.com/facebookgo/clock"
)

// Store performs read-through cache operations
// This is not a production cache,
// but used to show tiered storage patterns in Go
type Store struct {
	dogStore *dogStore
	ttl      time.Duration
	clock    clock.Clock
}

// StoreConfig is a config for the store constructor
type StoreConfig struct {
	TTL           time.Duration
	DogCacheStore DogCacheStore
	DogReadStore  DogReadStore
}

// NewStore is a constructor for a store
func NewStore(config StoreConfig) *Store {
	return &Store{
		&dogStore{
			readStore:  config.DogReadStore,
			cacheStore: config.DogCacheStore,
		},
		config.TTL,
		clock.New(),
	}
}
