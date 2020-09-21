package graph

import (
	"context"
	"github.com/facebookgo/clock"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"bin/bork/pkg/models"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any
// dependencies you require here.

type FetchDog func(ctx context.Context, id uuid.UUID) (*models.Dog, error)
type CreateDog func(ctx context.Context, dog *models.Dog) (*models.Dog, error)
type UpdateDog func(ctx context.Context, dog *models.Dog) (*models.Dog, error)
type FetchDogs func(ctx context.Context) (*models.Dogs, error)

type Resolver struct {
	Clock       clock.Clock
	Logger      *zap.Logger
	FetchDbDog  FetchDog
	CreateDbDog CreateDog
	UpdateDbDog UpdateDog
	FetchDbDogs FetchDogs
}
