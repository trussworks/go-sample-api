package postgres

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"bin/bork/pkg/apperrors"
	"bin/bork/pkg/models"
)

func TestFetchDog(t *testing.T) {
	store, clock, ctx, refresh := setup(t)

	t.Run("fetches dog when exists", func(t *testing.T) {
		refresh(t)
		insertDog := models.Dog{
			ID:        uuid.New(),
			Name:      "Chihua",
			Breed:     models.Chihuahua,
			BirthDate: clock.Now(),
			OwnerID:   "Owner",
		}
		expectedDog, err := store.CreateDog(ctx, &insertDog)
		assert.NoError(t, err)

		dog, err := store.FetchDog(ctx, expectedDog.ID)

		assert.NoError(t, err)
		assert.Equal(t, expectedDog.ID, dog.ID)
		assert.Equal(t, expectedDog.Name, dog.Name)
		assert.Equal(t, expectedDog.Breed, dog.Breed)
		assert.True(t, dog.BirthDate.Equal(expectedDog.BirthDate))
	})

	t.Run("returns error when doesn't exist", func(t *testing.T) {
		refresh(t)
		dog, err := store.FetchDog(ctx, uuid.New())

		assert.Error(t, err)
		assert.IsType(t, &apperrors.ResourceNotFoundError{}, err)
		assert.Nil(t, dog)
	})
}

func (s StoreTestSuite) TestCreateDog() {
	ctx := context.Background()

	s.Run("returns dog on success", func() {
		dog := models.Dog{
			ID:        uuid.UUID{},
			Name:      "Chihua",
			Breed:     models.Chihuahua,
			BirthDate: s.clock.Now(),
			OwnerID:   "Owner",
		}

		actualDog, err := s.store.CreateDog(ctx, &dog)

		s.NoError(err)
		s.NotZero(actualDog.ID)
	})

	conflictTests := []struct {
		dog     models.Dog
		message string
	}{
		{
			dog: models.Dog{
				ID:   uuid.UUID{},
				Name: "",
			},
			message: "pq: invalid input value for enum dog_breed: \"\"",
		},
		{
			dog: models.Dog{
				ID:    uuid.UUID{},
				Name:  "",
				Breed: models.Chihuahua,
			},
			message: "pq: new row for relation \"dog\" violates check constraint \"dog_name_check\"",
		},
		{
			dog: models.Dog{
				ID:    uuid.UUID{},
				Name:  "Lola",
				Breed: models.Chihuahua,
			},
			message: "pq: new row for relation \"dog\" violates check constraint \"dog_owner_id_check\"",
		},
	}
	for _, v := range conflictTests {
		s.Run("returns errors conflict failure", func() {
			actualDog, err := s.store.CreateDog(ctx, &v.dog)

			s.Error(err)
			s.Equal(v.message, err.Error())
			s.Nil(actualDog)
		})
	}
}

func (s StoreTestSuite) TestUpdateDog() {
	ctx := context.Background()

	s.Run("returns dog on success", func() {
		dog := models.Dog{
			ID:        uuid.UUID{},
			Name:      "Lola",
			Breed:     models.Chihuahua,
			BirthDate: s.clock.Now(),
			OwnerID:   "Owner",
		}
		createdDog, err := s.store.CreateDog(ctx, &dog)
		s.NoError(err)
		createdDog.Name = "Lolita"

		actualDog, err := s.store.UpdateDog(ctx, &dog)

		s.NoError(err)
		s.Equal("Lolita", actualDog.Name)
	})

	s.Run("fails if dog doesn't exist", func() {
		dog := models.Dog{
			ID:        uuid.New(),
			Name:      "Chihua",
			Breed:     models.Chihuahua,
			BirthDate: s.clock.Now(),
			OwnerID:   "Owner",
		}

		actualDog, err := s.store.UpdateDog(ctx, &dog)

		s.Error(err)
		s.IsType(&apperrors.ResourceNotFoundError{}, err)
		s.Nil(actualDog)
	})
}

func (s StoreTestSuite) TestFetchDogs() {
	ctx := context.Background()

	s.Run("fetches dogs", func() {
		insertDog := models.Dog{
			ID:        uuid.New(),
			Name:      "Chihua",
			Breed:     models.Chihuahua,
			BirthDate: s.clock.Now(),
			OwnerID:   "Owner",
		}
		expectedDog, err := s.store.CreateDog(ctx, &insertDog)
		s.NoError(err)

		dogs, err := s.store.FetchDogs(ctx)

		s.NoError(err)
		s.Len(*dogs, 1)
		s.Equal(expectedDog.ID, (*dogs)[0].ID)
	})
}
