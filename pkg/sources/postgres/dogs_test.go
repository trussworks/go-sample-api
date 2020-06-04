package postgres

import (
	"github.com/google/uuid"

	"bin/bork/pkg/models"
)

func (s StoreTestSuite) TestFetchDog() {
	const insertSQL = `
		INSERT INTO dog (id, name, breed, birth_date) VALUES (:id, :name, :breed, :birth_date)
	`
	s.Run("fetches dog when exists", func() {
		expectedDog := models.Dog{
			ID:        uuid.New(),
			Name:      "Chihua",
			Breed:     models.Chihuahua,
			BirthDate: s.clock.Now(),
		}
		_, err := s.db.NamedExec(insertSQL, &expectedDog)
		s.NoError(err)

		dog, err := s.store.FetchDog(expectedDog.ID)

		s.NoError(err)
		s.Equal(expectedDog.ID, dog.ID)
		s.Equal(expectedDog.Name, dog.Name)
		s.Equal(expectedDog.Breed, dog.Breed)
		s.True(dog.BirthDate.Equal(expectedDog.BirthDate))
	})
}
