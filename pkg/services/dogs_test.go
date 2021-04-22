package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"bin/bork/pkg/appcontext"
	"bin/bork/pkg/apperrors"
	"bin/bork/pkg/models"
)

func (s ServicesTestSuite) NewDog() models.Dog {
	return models.Dog{
		ID:        uuid.New(),
		Name:      "Lola",
		Breed:     models.Chihuahua,
		BirthDate: s.ServiceFactory.clock.Now(),
	}
}

func (s ServicesTestSuite) TestNewAuthorizeFetchDog() {
	authorize := NewAuthorizeFetchDog()
	s.Run("matching IDs returns true", func() {
		dog := models.Dog{
			OwnerID: "owner",
		}
		user := models.User{
			ID: dog.OwnerID,
		}

		ok, _ := authorize(user, &dog)

		s.True(ok)
	})

	s.Run("non-matching IDs returns true", func() {
		dog := models.Dog{
			OwnerID: "owner",
		}
		user := models.User{
			ID: "other owner",
		}

		ok, _ := authorize(user, &dog)

		s.False(ok)
	})
}

func (s ServicesTestSuite) TestServiceFactory_NewFetchDog() {
	expectedDog := models.Dog{
		ID:        uuid.New(),
		Name:      "Lola",
		Breed:     models.Chihuahua,
		BirthDate: s.ServiceFactory.clock.Now(),
	}
	fetch := func(ctx context.Context, uuid uuid.UUID) (*models.Dog, error) {
		return &expectedDog, nil
	}

	authorize := func(user models.User, dog *models.Dog) (bool, error) {
		return true, nil
	}

	s.Run("returns dog on golden path", func() {
		fetchDog := s.ServiceFactory.NewFetchDog(
			authorize,
			fetch,
		)
		ctx := appcontext.WithEmptyRequestLog(context.Background())
		ctx = appcontext.WithUser(ctx, models.User{})

		dog, err := fetchDog(ctx, uuid.New())

		s.NoError(err)
		s.Equal(&expectedDog, dog)
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
		ctx := appcontext.WithEmptyRequestLog(context.Background())
		ctx = appcontext.WithUser(ctx, models.User{})

		dog, err := fetchDog(ctx, uuid.New())

		s.IsType(&apperrors.UnauthorizedError{}, err)
		s.Nil(dog)
	})

	s.Run("logs the user id", func() {
		noAuthorize := func(models.User, *models.Dog) (bool, error) {
			return false, nil
		}
		fetchDog := s.ServiceFactory.NewFetchDog(
			noAuthorize,
			fetch,
		)
		ctx := appcontext.WithEmptyRequestLog(context.Background())
		ctx = appcontext.WithUser(ctx, models.User{
			ID: "foo",
		})

		_, _ = fetchDog(ctx, uuid.New())

		fields, ok := appcontext.RequestLogFields(ctx)
		if !ok {
			s.T().Fatal("couldn't get the fields")
		}

		s.Equal(zap.String("user_id", "foo"), fields[0])

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
		ctx := appcontext.WithEmptyRequestLog(context.Background())
		ctx = appcontext.WithUser(ctx, models.User{})

		dog, err := fetchDog(ctx, uuid.New())

		s.True(errors.Is(err, authErr))
		s.Nil(dog)
	})

	s.Run("returns error when fetch returns error", func() {
		fetchErr := errors.New("failed to fetch")
		failFetch := func(ctx context.Context, id uuid.UUID) (*models.Dog, error) {
			return &expectedDog, fetchErr
		}
		fetchDog := s.ServiceFactory.NewFetchDog(
			authorize,
			failFetch,
		)
		ctx := appcontext.WithEmptyRequestLog(context.Background())
		ctx = appcontext.WithUser(ctx, models.User{})

		dog, err := fetchDog(ctx, uuid.New())

		s.IsType(&apperrors.QueryError{}, err)
		s.Nil(dog)
	})
}

func (s ServicesTestSuite) TestNewAuthorizeCreateDog() {
	authorize := NewAuthorizeCreateDog()
	s.Run("any ID returns true", func() {
		dog := models.Dog{}
		user := models.User{
			ID: "owner",
		}

		ok, _ := authorize(user, &dog)

		s.True(ok)
	})

	s.Run("empty ID returns false", func() {
		dog := models.Dog{}
		user := models.User{
			ID: "",
		}

		ok, _ := authorize(user, &dog)

		s.False(ok)
	})
}

func (s ServicesTestSuite) TestServiceFactory_NewCreateDog() {
	create := func(ctx context.Context, dog *models.Dog) (*models.Dog, error) {
		return dog, nil
	}

	authorize := func(user models.User, dog *models.Dog) (bool, error) {
		return true, nil
	}

	s.Run("returns dog on golden path", func() {
		dog := s.NewDog()
		createDog := s.ServiceFactory.NewCreateDog(
			authorize,
			create,
		)
		ctx := appcontext.WithEmptyRequestLog(context.Background())
		ctx = appcontext.WithUser(ctx, models.User{})

		actualDog, err := createDog(ctx, &dog)

		s.NoError(err)
		s.NotZero(actualDog.ID)
		s.Equal(dog.Name, actualDog.Name)
	})

	s.Run("returns error with no user context", func() {
		dog := s.NewDog()
		createDog := s.ServiceFactory.NewCreateDog(
			authorize,
			create,
		)
		ctx := context.Background()

		actualDog, err := createDog(ctx, &dog)

		s.IsType(&apperrors.ContextError{}, err)
		s.Nil(actualDog)
	})

	s.Run("returns error when not authorized", func() {
		dog := s.NewDog()
		noAuthorize := func(models.User, *models.Dog) (bool, error) {
			return false, nil
		}
		createDog := s.ServiceFactory.NewCreateDog(
			noAuthorize,
			create,
		)
		ctx := appcontext.WithEmptyRequestLog(context.Background())
		ctx = appcontext.WithUser(ctx, models.User{})

		actualDog, err := createDog(ctx, &dog)

		s.IsType(&apperrors.UnauthorizedError{}, err)
		s.Nil(actualDog)
	})

	s.Run("returns error when authorize returns error", func() {
		dog := s.NewDog()
		authErr := errors.New("failed to authorize")
		failAuthorize := func(models.User, *models.Dog) (bool, error) {
			return false, authErr
		}
		createDog := s.ServiceFactory.NewCreateDog(
			failAuthorize,
			create,
		)
		ctx := appcontext.WithEmptyRequestLog(context.Background())
		ctx = appcontext.WithUser(ctx, models.User{})

		actualDog, err := createDog(ctx, &dog)

		s.Equal(authErr, err)
		s.Nil(actualDog)
	})

	s.Run("returns error when create returns error", func() {
		createdDog := s.NewDog()
		fetchErr := errors.New("failed to create")
		failCreate := func(ctx context.Context, dog *models.Dog) (*models.Dog, error) {
			return &createdDog, fetchErr
		}
		createDog := s.ServiceFactory.NewCreateDog(
			authorize,
			failCreate,
		)
		ctx := appcontext.WithEmptyRequestLog(context.Background())
		ctx = appcontext.WithUser(ctx, models.User{})

		actualDog, err := createDog(ctx, &createdDog)

		s.IsType(&apperrors.QueryError{}, err)
		s.Nil(actualDog)
	})
}

func (s ServicesTestSuite) TestNewAuthorizeUpdateDog() {
	authorize := NewAuthorizeUpdateDog()
	s.Run("matching IDs returns true", func() {
		dog := models.Dog{
			OwnerID: "owner",
		}
		user := models.User{
			ID: dog.OwnerID,
		}

		ok, _ := authorize(user, &dog)

		s.True(ok)
	})

	s.Run("non matching IDs returns false", func() {
		dog := models.Dog{
			OwnerID: "owner",
		}
		user := models.User{
			ID: "other owner",
		}

		ok, _ := authorize(user, &dog)

		s.False(ok)
	})
}

func (s ServicesTestSuite) TestServiceFactory_NewUpdateDog() {
	update := func(ctx context.Context, dog *models.Dog) (*models.Dog, error) {
		return dog, nil
	}

	fetch := func(ctx context.Context, uuid uuid.UUID) (*models.Dog, error) {
		return &models.Dog{}, nil
	}

	authorize := func(user models.User, dog *models.Dog) (bool, error) {
		return true, nil
	}

	s.Run("returns dog on golden path", func() {
		dog := s.NewDog()
		updateDog := s.ServiceFactory.NewUpdateDog(
			authorize,
			update,
			fetch,
		)
		ctx := appcontext.WithEmptyRequestLog(context.Background())
		ctx = appcontext.WithUser(ctx, models.User{})

		actualDog, err := updateDog(ctx, &dog)

		s.NoError(err)
		s.NotZero(actualDog.ID)
		s.Equal(dog.Name, actualDog.Name)
	})

	s.Run("returns error with no user context", func() {
		dog := s.NewDog()
		updateDog := s.ServiceFactory.NewUpdateDog(
			authorize,
			update,
			fetch,
		)
		ctx := context.Background()

		actualDog, err := updateDog(ctx, &dog)

		s.IsType(&apperrors.ContextError{}, err)
		s.Nil(actualDog)
	})

	s.Run("returns error when not authorized", func() {
		dog := s.NewDog()
		noAuthorize := func(models.User, *models.Dog) (bool, error) {
			return false, nil
		}
		updateDog := s.ServiceFactory.NewUpdateDog(
			noAuthorize,
			update,
			fetch,
		)
		ctx := appcontext.WithEmptyRequestLog(context.Background())
		ctx = appcontext.WithUser(ctx, models.User{})

		actualDog, err := updateDog(ctx, &dog)

		s.IsType(&apperrors.UnauthorizedError{}, err)
		s.Nil(actualDog)
	})

	s.Run("returns error when authorize returns error", func() {
		dog := s.NewDog()
		authErr := errors.New("failed to authorize")
		failAuthorize := func(models.User, *models.Dog) (bool, error) {
			return false, authErr
		}
		updateDog := s.ServiceFactory.NewUpdateDog(
			failAuthorize,
			update,
			fetch,
		)
		ctx := appcontext.WithEmptyRequestLog(context.Background())
		ctx = appcontext.WithUser(ctx, models.User{})

		actualDog, err := updateDog(ctx, &dog)

		s.Equal(authErr, err)
		s.Nil(actualDog)
	})

	s.Run("returns error when fetch fail", func() {
		dog := s.NewDog()
		failFetch := func(ctx context.Context, id uuid.UUID) (*models.Dog, error) {
			return nil, errors.New("failed to fetch")
		}
		updateDog := s.ServiceFactory.NewUpdateDog(
			authorize,
			update,
			failFetch,
		)
		ctx := appcontext.WithEmptyRequestLog(context.Background())
		ctx = appcontext.WithUser(ctx, models.User{})

		actualDog, err := updateDog(ctx, &dog)

		s.IsType(&apperrors.QueryError{}, err)
		s.Nil(actualDog)
	})

	s.Run("returns error when update returns error", func() {
		updatedDog := s.NewDog()
		fetchErr := errors.New("failed to update")
		failUpdate := func(ctx context.Context, dog *models.Dog) (*models.Dog, error) {
			return &updatedDog, fetchErr
		}
		updateDog := s.ServiceFactory.NewUpdateDog(
			authorize,
			failUpdate,
			fetch,
		)
		ctx := appcontext.WithEmptyRequestLog(context.Background())
		ctx = appcontext.WithUser(ctx, models.User{})

		actualDog, err := updateDog(ctx, &updatedDog)

		s.IsType(&apperrors.QueryError{}, err)
		s.Nil(actualDog)
	})
}

func (s ServicesTestSuite) TestNewAuthorizeFetchDogs() {
	authorize := NewAuthorizeFetchDogs()
	s.Run("any ID returns true", func() {
		user := models.User{
			ID: "owner",
		}

		ok, _ := authorize(user)

		s.True(ok)
	})

	s.Run("empty ID returns false", func() {
		user := models.User{
			ID: "",
		}

		ok, _ := authorize(user)

		s.False(ok)
	})
}

func (s ServicesTestSuite) TestServiceFactory_NewFetchDogs() {
	expectedDog := models.Dog{
		ID:        uuid.New(),
		Name:      "Lola",
		Breed:     models.Chihuahua,
		BirthDate: s.ServiceFactory.clock.Now(),
	}
	fetch := func(ctx context.Context) (*models.Dogs, error) {
		return &models.Dogs{expectedDog}, nil
	}

	authorize := func(user models.User) (bool, error) {
		return true, nil
	}

	s.Run("returns dogs on golden path", func() {
		fetchDogs := s.ServiceFactory.NewFetchDogs(
			authorize,
			fetch,
		)
		ctx := appcontext.WithEmptyRequestLog(context.Background())
		ctx = appcontext.WithUser(ctx, models.User{})

		actualDogs, err := fetchDogs(ctx)

		s.NoError(err)
		s.Equal(models.Dogs{expectedDog}, *actualDogs)
	})

	s.Run("returns error with no user context", func() {
		fetchDogs := s.ServiceFactory.NewFetchDogs(
			authorize,
			fetch,
		)
		ctx := context.Background()

		dog, err := fetchDogs(ctx)

		s.IsType(&apperrors.ContextError{}, err)
		s.Nil(dog)
	})

	s.Run("returns error when not authorized", func() {
		noAuthorize := func(models.User) (bool, error) {
			return false, nil
		}
		fetchDogs := s.ServiceFactory.NewFetchDogs(
			noAuthorize,
			fetch,
		)
		ctx := appcontext.WithEmptyRequestLog(context.Background())
		ctx = appcontext.WithUser(ctx, models.User{})

		dog, err := fetchDogs(ctx)

		s.IsType(&apperrors.UnauthorizedError{}, err)
		s.Nil(dog)
	})

	s.Run("returns error when authorize returns error", func() {
		authErr := errors.New("failed to authorize")
		failAuthorize := func(models.User) (bool, error) {
			return false, authErr
		}
		fetchDogs := s.ServiceFactory.NewFetchDogs(
			failAuthorize,
			fetch,
		)
		ctx := appcontext.WithEmptyRequestLog(context.Background())
		ctx = appcontext.WithUser(ctx, models.User{})

		dog, err := fetchDogs(ctx)

		s.Equal(authErr, err)
		s.Nil(dog)
	})

	s.Run("returns error when fetch returns error", func() {
		fetchErr := errors.New("failed to fetch")
		failFetch := func(ctx context.Context) (*models.Dogs, error) {
			return &models.Dogs{expectedDog}, fetchErr
		}
		fetchDogs := s.ServiceFactory.NewFetchDogs(
			authorize,
			failFetch,
		)
		ctx := appcontext.WithEmptyRequestLog(context.Background())
		ctx = appcontext.WithUser(ctx, models.User{})

		dog, err := fetchDogs(ctx)

		s.IsType(&apperrors.QueryError{}, err)
		s.Nil(dog)
	})
}
