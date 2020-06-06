package cache

import (
	"time"

	"bin/bork/pkg/models"
)

type DogCacheStore interface {
	FetchDogs() (*models.Dogs, error)
	SaveDogs(*models.Dogs) error
}

type DogReadStore interface {
	FetchDogs() (*models.Dogs, error)
}

type dogStore struct {
	updatedAt  time.Time
	readStore  DogReadStore
	cacheStore DogCacheStore
}

func (s *Store) FetchDogs() (*models.Dogs, error) {
	if s.clock.Now().Before(s.dogStore.updatedAt.Add(s.ttl)) {
		return s.dogStore.cacheStore.FetchDogs()
	}
	dogs, err := s.dogStore.readStore.FetchDogs()
	if err != nil {
		return nil, err
	}
	err = s.dogStore.cacheStore.SaveDogs(dogs)
	if err != nil {
		return dogs, err
	}
	s.dogStore.updatedAt = s.clock.Now()
	return dogs, nil
}
