package cache

import (
	"context"

	"github.com/google/uuid"

	"bin/bork/pkg/models"
)

type fakeDogReadStore struct {
	Dogs models.Dogs
}

func (s fakeDogReadStore) FetchDogs(context.Context) (*models.Dogs, error) {
	return &s.Dogs, nil
}

func (s StoreTestSuite) TestFetchDogs() {
	readStore := fakeDogReadStore{models.Dogs{
		models.Dog{
			ID:        uuid.New(),
			Name:      "Chihua",
			Breed:     models.Chihuahua,
			BirthDate: s.clock.Now(),
			OwnerID:   "Owner",
		},
	},
	}
	s.store.dogStore.readStore = readStore
	ctx := context.Background()

	s.Run("first fetch gets read store dogs", func() {
		actualDogs, err := s.store.FetchDogs(ctx)

		s.NoError(err)
		s.Equal(readStore.Dogs, *actualDogs)
		s.True(s.store.dogStore.updatedAt.Equal(s.clock.Now()))
	})

	s.Run("under TTL gets cache store dogs", func() {
		oldTime := s.clock.Now()
		s.clock.Add(s.store.ttl / 2)
		oldDogs := readStore.Dogs
		readStore.Dogs = append(readStore.Dogs, models.Dog{
			ID:        uuid.New(),
			Name:      "Lola",
			Breed:     models.Chihuahua,
			BirthDate: s.clock.Now(),
			OwnerID:   "Owner",
		})
		s.store.dogStore.readStore = readStore

		actualDogs, err := s.store.FetchDogs(ctx)

		s.NoError(err)
		s.Equal(oldDogs, *actualDogs)
		s.True(s.store.dogStore.updatedAt.Equal(oldTime))
	})

	s.Run("over TTL gets read store dogs", func() {
		s.clock.Add(s.store.ttl)

		actualDogs, err := s.store.FetchDogs(ctx)

		s.NoError(err)
		s.Equal(readStore.Dogs, *actualDogs)
		s.True(s.store.dogStore.updatedAt.Equal(s.clock.Now()))
	})
}
