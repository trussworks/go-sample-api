package memory

import (
	"context"

	"github.com/google/uuid"

	"bin/bork/pkg/models"
)

func (s StoreTestSuite) TestFetchDogs() {
	s.store.dogs = &models.Dogs{
		models.Dog{
			ID:        uuid.New(),
			Name:      "Chihua",
			Breed:     models.Chihuahua,
			BirthDate: s.clock.Now(),
			OwnerID:   "Owner",
		},
	}

	dogs, err := s.store.FetchDogs(context.Background())

	s.NoError(err)
	s.Equal(s.store.dogs, dogs)
}

func (s StoreTestSuite) TestSaveDogs() {
	expectedDogs := &models.Dogs{
		models.Dog{
			ID:        uuid.New(),
			Name:      "Chihua",
			Breed:     models.Chihuahua,
			BirthDate: s.clock.Now(),
			OwnerID:   "Owner",
		},
	}

	err := s.store.SaveDogs(context.Background(), expectedDogs)

	s.NoError(err)
	s.Equal(expectedDogs, s.store.dogs)
}
