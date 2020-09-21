package model

import (
	"fmt"

	"bin/bork/pkg/models"
)

func GqlDogInputToDbDog(newDog DogInput) (*models.Dog, error) {

	d := &models.Dog{}
	switch (newDog.Breed) {
	case BreedChihuahua:
		d.Breed = models.Chihuahua
	default:
		return nil, fmt.Errorf("%s is not a valid Breed", newDog.Breed)
	}
	d.Name = newDog.Name
	d.BirthDate = newDog.BirthDate
	return d, nil
}

func DbDogToGqlDog(dog *models.Dog) (*Dog, error)  {
	d := &Dog{}
	switch(dog.Breed) {
	case models.Chihuahua:
		d.Breed = BreedChihuahua
	default:
		return nil, fmt.Errorf("%s is not a valid Breed", dog.Breed)
	}
	d.BirthDate = dog.BirthDate
	d.ID = dog.ID.String()
	d.Name = dog.Name
	d.Owner = &Owner{
		ID: dog.OwnerID,
	}
	return d, nil
}
