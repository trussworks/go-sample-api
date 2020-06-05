package memory

import (
	"testing"

	"github.com/facebookgo/clock"
	"github.com/stretchr/testify/suite"
)

type StoreTestSuite struct {
	suite.Suite
	store *Store
	clock *clock.Mock
}

func TestStoreTestSuite(t *testing.T) {
	store := NewStore()

	storeTestSuite := &StoreTestSuite{
		Suite: suite.Suite{},
		store: store,
		clock: clock.NewMock(),
	}

	suite.Run(t, storeTestSuite)
}
