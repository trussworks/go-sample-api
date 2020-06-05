package memory

import "bin/bork/pkg/models"

// Store performs memory db operations
// This is not a production store,
// but used to show variance in stores
type Store struct {
	dogs *models.Dogs
}

// NewStore is a constructor for a store
func NewStore() *Store {
	return &Store{
		&models.Dogs{},
	}
}
