package postgres

import (
	"github.com/google/uuid"

	"bin/bork/pkg/models"
)

// FetchDog queries the DB for a dog
func (s *Store) FetchDog(id uuid.UUID) (*models.Dog, error) {
	dog := models.Dog{}
	const fetchDogSQL = `
		SELECT
			*
		FROM
			dog
		WHERE
			dog.id = $1`

	err := s.db.Get(&dog, fetchDogSQL, id)
	if err != nil {
		return &models.Dog{}, err
	}
	return &dog, nil
}
