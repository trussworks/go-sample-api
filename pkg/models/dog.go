package models

import (
	"time"

	"github.com/google/uuid"
)

// DogBreed is a type for dog breeds
type DogBreed string

const (
	Chihuahua DogBreed = "Chihuahua"
)

// Dog is a model for dog fields
type Dog struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Breed     DogBreed  `json:"breed" db:"breed"`
	BirthDate time.Time `json:"birthDate" db:"birth_date"`
	OwnerID   string    `json:"owner_id" db:"owner_id"`
}
