package postgres

import (
	"fmt"
	"testing"

	"github.com/facebookgo/clock"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // required for postgres driver in sqlx
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"

	"bin/bork/pkg/appconfig"
)

type StoreTestSuite struct {
	suite.Suite
	db    *sqlx.DB
	store *Store
	clock *clock.Mock
}

func TestStoreTestSuite(t *testing.T) {
	config := viper.New()
	config.AutomaticEnv()

	dbConfig := DBConfig{
		Host:     config.GetString(appconfig.DBHostConfigKey),
		Port:     config.GetString(appconfig.DBPortConfigKey),
		Database: config.GetString(appconfig.DBNameConfigKey),
		Username: config.GetString(appconfig.DBUsernameConfigKey),
		Password: config.GetString(appconfig.DBPasswordConfigKey),
		SSLMode:  config.GetString(appconfig.DBSSLModeConfigKey),
	}
	store, err := NewStore(dbConfig)
	if err != nil {
		fmt.Printf("Failed to get new database: %v", err)
		t.Fail()
	}

	storeTestSuite := &StoreTestSuite{
		Suite: suite.Suite{},
		db:    store.db,
		store: store,
		clock: clock.NewMock(),
	}

	suite.Run(t, storeTestSuite)
}
