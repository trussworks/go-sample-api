package cache

import (
	"context"
	"time"

	"bin/bork/pkg/models"
)

type DogCacheStore interface {
	FetchDogs(context.Context) (*models.Dogs, error)
	SaveDogs(context.Context, *models.Dogs) error
}

type DogReadStore interface {
	FetchDogs(ctx context.Context) (*models.Dogs, error)
}

type dogStore struct {
	updatedAt  time.Time
	readStore  DogReadStore
	cacheStore DogCacheStore
}

func (s *Store) FetchDogs(ctx context.Context) (*models.Dogs, error) {
	if s.clock.Now().Before(s.dogStore.updatedAt.Add(s.ttl)) {
		return s.dogStore.cacheStore.FetchDogs(ctx)
	}
	dogs, err := s.dogStore.readStore.FetchDogs(ctx)
	if err != nil {
		return nil, err
	}
	err = s.dogStore.cacheStore.SaveDogs(ctx, dogs)
	if err != nil {
		return dogs, err
	}
	s.dogStore.updatedAt = s.clock.Now()
	return dogs, nil
}
