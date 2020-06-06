package cache

import (
	"testing"
	"time"

	"github.com/facebookgo/clock"
	"github.com/stretchr/testify/suite"

	"bin/bork/pkg/sources/memory"
)

type StoreTestSuite struct {
	suite.Suite
	store *Store
	clock *clock.Mock
}

func TestStoreTestSuite(t *testing.T) {
	config := StoreConfig{
		TTL:           time.Minute,
		DogCacheStore: memory.NewStore(),
		DogReadStore:  fakeDogReadStore{},
	}
	store := NewStore(config)
	mockClock := clock.NewMock()
	store.clock = mockClock

	storeTestSuite := &StoreTestSuite{
		Suite: suite.Suite{},
		store: store,
		clock: mockClock,
	}

	suite.Run(t, storeTestSuite)
}
