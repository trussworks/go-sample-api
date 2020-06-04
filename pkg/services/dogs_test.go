package services

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"bin/bork/pkg/appcontext"
	"bin/bork/pkg/apperrors"
	"bin/bork/pkg/models"
)

func (s ServicesTestSuite) TestNewAuthorizeFetchDog() {
	authorize := NewAuthorizeFetchDog()
	s.Run("matching IDs returns true", func() {
		dog := models.Dog{
			OwnerID:   "owner",
		}
		user := models.User{
			ID:    dog.OwnerID,
		}

		ok, _ := authorize(user, &dog)

		s.True(ok)
	})

	s.Run("non-matching IDs returns true", func() {
		dog := models.Dog{
			OwnerID:   "owner",
		}
		user := models.User{
			ID:    "other owner",
		}

		ok, _ := authorize(user, &dog)

		s.False(ok)
	})
}

func (s ServicesTestSuite) TestServiceFactory_NewFetchDog() {
	fetchedDog := models.Dog{
		ID:        uuid.New(),
		Name:      "Lola",
		Breed:     models.Chihuahua,
		BirthDate: s.ServiceFactory.clock.Now(),
	}
	fetch := func(uuid uuid.UUID) (*models.Dog, error) {
		return &fetchedDog, nil
	}

	authorize := func(user models.User, dog *models.Dog) (bool, error) {
		return true, nil
	}

	s.Run("returns dog on golden path", func() {
		fetchDog := s.ServiceFactory.NewFetchDog(
			authorize,
			fetch,
		)
		ctx := context.Background()
		ctx = appcontext.WithUser(ctx, models.User{})

		dog, err := fetchDog(ctx, uuid.New())

		s.NoError(err)
		s.Equal(&fetchedDog, dog)
	})

	s.Run("returns error with no user context", func() {
		fetchDog := s.ServiceFactory.NewFetchDog(
			authorize,
			fetch,
		)
		ctx := context.Background()

		dog, err := fetchDog(ctx, uuid.New())

		s.IsType(&apperrors.ContextError{}, err)
		s.Nil(dog)
	})

	s.Run("returns error when not authorized", func() {
		noAuthorize := func(models.User, *models.Dog) (bool, error) {
			return false, nil
		}
		fetchDog := s.ServiceFactory.NewFetchDog(
			noAuthorize,
			fetch,
		)
		ctx := context.Background()
		ctx = appcontext.WithUser(ctx, models.User{})

		dog, err := fetchDog(ctx, uuid.New())

		s.IsType(&apperrors.UnauthorizedError{}, err)
		s.Nil(dog)
	})

	s.Run("returns error when authorize returns error", func() {
		authErr := errors.New("failed to authorize")
		failAuthorize := func(models.User, *models.Dog) (bool, error) {
			return false, authErr
		}
		fetchDog := s.ServiceFactory.NewFetchDog(
			failAuthorize,
			fetch,
		)
		ctx := context.Background()
		ctx = appcontext.WithUser(ctx, models.User{})

		dog, err := fetchDog(ctx, uuid.New())

		s.Equal(authErr, err)
		s.Nil(dog)
	})

	s.Run("returns error when fetch returns error", func() {
		fetchErr := errors.New("failed to fetch")
		failFetch := func(id uuid.UUID) (*models.Dog, error) {
			return &fetchedDog, fetchErr
		}
		fetchDog := s.ServiceFactory.NewFetchDog(
			authorize,
			failFetch,
		)
		ctx := context.Background()
		ctx = appcontext.WithUser(ctx, models.User{})

		dog, err := fetchDog(ctx, uuid.New())

		s.IsType(&apperrors.QueryError{}, err)
		s.Nil(dog)
	})
}
