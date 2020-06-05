package memory

import "bin/bork/pkg/models"

// FetchDogs returns dogs in the memory store
func (s *Store) FetchDogs() (*models.Dogs, error) {
	return s.dogs, nil
}

// SaveDogs saves dogs to the memory store
func (s *Store) SaveDogs(dogs *models.Dogs) error {
	s.dogs = dogs
	return nil
}
