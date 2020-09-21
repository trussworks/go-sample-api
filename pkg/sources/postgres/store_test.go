package postgres

import (
	"testing"

	"github.com/facebookgo/clock"
	_ "github.com/lib/pq" // required for postgres driver in sqlx
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"

	"bin/bork/pkg/appconfig"
)

type StoreTestSuite struct {
	suite.Suite
	store *Store
	clock *clock.Mock
}

func TestStoreTestSuite(t *testing.T) {
	config := viper.New()
	config.AutomaticEnv()

	origAppConfig, err := appconfig.NewAppConfig(config)
	if err != nil {
		t.Fatalf("Error initializing config: %s", err)
	}

	store, err := NewStore(&origAppConfig)
	if err != nil {
		t.Fatalf("Error initializing store: %s", err)
	}

	storeTestSuite := &StoreTestSuite{
		Suite: suite.Suite{},
		clock: clock.NewMock(),
		store: store,
	}

	suite.Run(t, storeTestSuite)
}

func (s *StoreTestSuite) SetupTest() {
	s.store.db.MustExec("TRUNCATE dog")
}
