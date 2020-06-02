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
	ID        uuid.UUID
	Name      string
	Breed     DogBreed
	BirthDate time.Time
}
