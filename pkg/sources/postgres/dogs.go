package postgres

import (
	"github.com/google/uuid"

	"bin/bork/pkg/apperrors"
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
		if err.Error() == "sql: no rows in result set" {
			return nil, &apperrors.ResourceNotFoundError{Resource: dog}
		}
		return nil, err
	}
	return &dog, nil
}

// CreateDog creates a dog in the DB
func (s *Store) CreateDog(dog *models.Dog) (*models.Dog, error) {
	dog.ID = uuid.New()
	const createDogSQL = `
		INSERT INTO dog (
		                 id,
		                 name,
		                 breed,
		                 birth_date,
		                 owner_id
		)
		VALUES (
		        :id,
		        :name,
		        :breed,
		        :birth_date,
		        :owner_id
		)`

	_, err := s.db.NamedExec(createDogSQL, &dog)
	if err != nil {
		return nil, err
	}
	return dog, nil
}

// UpdateDog creates a dog in the DB
func (s *Store) UpdateDog(dog *models.Dog) (*models.Dog, error) {
	const updateDogSQL = `
		UPDATE dog 
		SET
		    name = :name,
		    breed = :breed,
		    birth_date = :birth_date
		WHERE dog.id = :id
		`

	result, err := s.db.NamedExec(updateDogSQL, &dog)
	if err != nil {
		return nil, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, &apperrors.ResourceNotFoundError{Resource: dog}
	}
	return dog, nil
}

// FetchDogs queries the DB for dogs
func (s *Store) FetchDogs() (*models.Dogs, error) {
	dog := models.Dogs{}
	const fetchDogsSQL = `
		SELECT
			*
		FROM
			dog`

	err := s.db.Select(&dog, fetchDogsSQL)
	if err != nil {
		return nil, err
	}
	return &dog, nil
}
