package graph

import (
	"github.com/facebookgo/clock"
	"go.uber.org/zap"

	"bin/bork/pkg/models"
	"bin/bork/pkg/sources/postgres"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any
// dependencies you require here.

type AuthorizeFetchDog func(user models.User, dog *models.Dog) (bool, error)
type AuthorizeCreateDog func(user models.User, dog *models.Dog) (bool, error)
type Resolver struct {
	Clock  clock.Clock
	Logger *zap.Logger
	Store  *postgres.Store
	AuthorizeFetchDog AuthorizeFetchDog
	AuthorizeCreateDog AuthorizeCreateDog
}
