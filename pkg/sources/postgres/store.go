package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // required for postgres driver in sqlx

	"bin/bork/pkg/appconfig"
)

// Store performs database operations for the bork application
type Store struct {
	db *sqlx.DB
}


func BuildDataSourceName(appConfig *appconfig.AppConfig) string {
	var sslmode = "disable"
	if appConfig.DBSSLMode == "enable" {
		sslmode = "enable"
	}
	return fmt.Sprintf(
		"host=%s port=%s user=%s "+
			"password=%s dbname=%s sslmode=%s",
		appConfig.DBHost,
		appConfig.DBPort,
		appConfig.DBUsername,
		appConfig.DBPassword,
		appConfig.DBName,
		sslmode,
	)
}

func NewDB(appConfig *appconfig.AppConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", BuildDataSourceName(appConfig))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewStore(appConfig *appconfig.AppConfig) (*Store, error) {
	db, err := NewDB(appConfig)
	if err != nil {
		return nil, err
	}
	return NewStoreWithDB(db), nil
}

// NewStore is a constructor for a store
func NewStoreWithDB(db *sqlx.DB) *Store {
	return &Store{db: db}
}
