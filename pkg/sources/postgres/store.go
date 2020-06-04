package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // required for postgres driver in sqlx
)

// Store performs database operations for the bork application
type Store struct {
	db *sqlx.DB
}

// DBConfig holds the configurations for a database connection
type DBConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
	SSLMode  string
}

// NewStore is a constructor for a store
func NewStore(
	config DBConfig,
) (*Store, error) {
	dataSourceName := fmt.Sprintf(
		"host=%s port=%s user=%s "+
			"password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.Username,
		config.Password,
		config.Database,
	)
	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}
