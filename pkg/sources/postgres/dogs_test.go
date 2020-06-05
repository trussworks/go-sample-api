package postgres

import (
	"github.com/google/uuid"

	"bin/bork/pkg/apperrors"
	"bin/bork/pkg/models"
)

func (s StoreTestSuite) TestFetchDog() {
	s.Run("fetches dog when exists", func() {
		insertDog := models.Dog{
			ID:        uuid.New(),
			Name:      "Chihua",
			Breed:     models.Chihuahua,
			BirthDate: s.clock.Now(),
			OwnerID:   "Owner",
		}
		expectedDog, err := s.store.CreateDog(&insertDog)
		s.NoError(err)

		dog, err := s.store.FetchDog(expectedDog.ID)

		s.NoError(err)
		s.Equal(expectedDog.ID, dog.ID)
		s.Equal(expectedDog.Name, dog.Name)
		s.Equal(expectedDog.Breed, dog.Breed)
		s.True(dog.BirthDate.Equal(expectedDog.BirthDate))
	})

	s.Run("returns error when doesn't exist", func() {
		dog, err := s.store.FetchDog(uuid.New())

		s.Error(err)
		s.IsType(&apperrors.ResourceNotFoundError{}, err)
		s.Nil(dog)
	})
}

func (s StoreTestSuite) TestCreateDog() {
	s.Run("returns dog on success", func() {
		dog := models.Dog{
			ID:        uuid.UUID{},
			Name:      "Chihua",
			Breed:     models.Chihuahua,
			BirthDate: s.clock.Now(),
			OwnerID:   "Owner",
		}

		actualDog, err := s.store.CreateDog(&dog)

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
			actualDog, err := s.store.CreateDog(&v.dog)

			s.Error(err)
			s.Equal(v.message, err.Error())
			s.Nil(actualDog)
		})
	}
}

func (s StoreTestSuite) TestUpdateDog() {
	s.Run("returns dog on success", func() {
		dog := models.Dog{
			ID:        uuid.UUID{},
			Name:      "Lola",
			Breed:     models.Chihuahua,
			BirthDate: s.clock.Now(),
			OwnerID:   "Owner",
		}
		createdDog, err := s.store.CreateDog(&dog)
		s.NoError(err)
		createdDog.Name = "Lolita"

		actualDog, err := s.store.UpdateDog(&dog)

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

		actualDog, err := s.store.UpdateDog(&dog)

		s.Error(err)
		s.IsType(&apperrors.ResourceNotFoundError{}, err)
		s.Nil(actualDog)
	})
}
