package memory

import (
	"bin/bork/pkg/models"
	"context"
)

// FetchDogs returns dogs in the memory store
func (s *Store) FetchDogs(ctx context.Context) (*models.Dogs, error) {
	return s.dogs, nil
}

// SaveDogs saves dogs to the memory store
func (s *Store) SaveDogs(ctx context.Context, dogs *models.Dogs) error {
	s.dogs = dogs
	return nil
}
