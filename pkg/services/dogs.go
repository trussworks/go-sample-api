package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"bin/bork/pkg/appcontext"
	"bin/bork/pkg/apperrors"
	"bin/bork/pkg/models"
)

// NewAuthorizeFetchDog authorizes fetching a dog
func NewAuthorizeFetchDog() func(user models.User, dog *models.Dog) (bool, error) {
	return func(user models.User, dog *models.Dog) (bool, error) {
		if dog.OwnerID == user.ID {
			return true, nil
		}
		return false, nil
	}
}

// NewFetchDog returns a service function for fetching a dog
func (f ServiceFactory) NewFetchDog(
	authorize func(user models.User, dog *models.Dog) (bool, error),
	fetch func(ctx context.Context, id uuid.UUID) (*models.Dog, error),
) func(ctx context.Context, id uuid.UUID) (*models.Dog, error) {
	return func(ctx context.Context, id uuid.UUID) (*models.Dog, error) {
		dog, err := fetch(ctx, id)
		if err != nil {
			queryError := apperrors.QueryError{
				Err:       err,
				Resource:  models.Dog{},
				Operation: apperrors.QueryFetch,
			}
			return nil, &queryError
		}
		user, ok := appcontext.User(ctx)
		if !ok {
			contextError := apperrors.ContextError{
				Err:       errors.New("failed to get context"),
				Resource:  apperrors.ContextResourceUser,
				Operation: apperrors.ContextOperationGet,
			}
			return nil, &contextError
		}
		appcontext.LogRequestField(ctx, zap.String("user_id", user.ID))

		ok, err = authorize(user, dog)
		if err != nil {
			return nil, fmt.Errorf("failed to authorize fetchDog: %w", err)
		}
		if !ok {
			unauthorizedError := apperrors.UnauthorizedError{
				User:      user,
				Operation: apperrors.QueryFetch,
				Resource:  dog,
				Err:       err,
			}
			return nil, &unauthorizedError
		}
		return dog, nil
	}
}

// NewAuthorizeCreateDog authorizing creating a dog
func NewAuthorizeCreateDog() func(user models.User, dog *models.Dog) (bool, error) {
	return func(user models.User, dog *models.Dog) (bool, error) {
		return user.ID != "", nil
	}
}

// NewCreateDog returns a service function for creating a dog
func (f ServiceFactory) NewCreateDog(
	authorize func(user models.User, dog *models.Dog) (bool, error),
	create func(ctx context.Context, dog *models.Dog) (*models.Dog, error),
) func(ctx context.Context, dog *models.Dog) (*models.Dog, error) {
	return func(ctx context.Context, dog *models.Dog) (*models.Dog, error) {
		user, ok := appcontext.User(ctx)
		if !ok {
			contextError := apperrors.ContextError{
				Err:       errors.New("failed to get context"),
				Resource:  apperrors.ContextResourceUser,
				Operation: apperrors.ContextOperationGet,
			}
			return nil, &contextError
		}
		appcontext.LogRequestField(ctx, zap.String("user", user.ID))
		ok, err := authorize(user, dog)
		if err != nil {
			appcontext.LogRequestError(ctx, "failed to authorize createDog", err)
			return nil, err
		}
		if !ok {
			unauthorizedError := apperrors.UnauthorizedError{
				User:      user,
				Operation: apperrors.QueryCreate,
				Resource:  dog,
				Err:       err,
			}
			return nil, &unauthorizedError
		}
		dog.OwnerID = user.ID
		createdDog, err := create(ctx, dog)
		if err != nil {
			queryError := apperrors.QueryError{
				Err:       err,
				Resource:  dog,
				Operation: apperrors.QueryCreate,
			}
			return nil, &queryError
		}
		return createdDog, nil
	}
}

// NewAuthorizeUpdateDog authorizing updating a dog
func NewAuthorizeUpdateDog() func(user models.User, dog *models.Dog) (bool, error) {
	return func(user models.User, dog *models.Dog) (bool, error) {
		if user.ID == dog.OwnerID {
			return true, nil
		}
		return false, nil
	}
}

// NewUpdateDog returns a service function for updating a dog
func (f ServiceFactory) NewUpdateDog(
	authorize func(user models.User, dog *models.Dog) (bool, error),
	update func(ctx context.Context, dog *models.Dog) (*models.Dog, error),
	fetch func(ctx context.Context, id uuid.UUID) (*models.Dog, error),
) func(ctx context.Context, dog *models.Dog) (*models.Dog, error) {
	return func(ctx context.Context, dog *models.Dog) (*models.Dog, error) {
		user, ok := appcontext.User(ctx)
		if !ok {
			contextError := apperrors.ContextError{
				Err:       errors.New("failed to get context"),
				Resource:  apperrors.ContextResourceUser,
				Operation: apperrors.ContextOperationGet,
			}
			return nil, &contextError
		}
		appcontext.LogRequestField(ctx, zap.String("user", user.ID))
		existingDog, err := fetch(ctx, dog.ID)
		if err != nil {
			queryError := apperrors.QueryError{
				Err:       err,
				Resource:  models.Dog{},
				Operation: apperrors.QueryUpdate,
			}
			return nil, &queryError
		}
		ok, err = authorize(user, existingDog)
		if err != nil {
			appcontext.LogRequestError(ctx, "failed to authorize updateDog", err)
			return nil, err
		}
		if !ok {
			unauthorizedError := apperrors.UnauthorizedError{
				User:      user,
				Operation: apperrors.QueryCreate,
				Resource:  dog,
				Err:       err,
			}
			return nil, &unauthorizedError
		}
		dog.OwnerID = user.ID
		createdDog, err := update(ctx, dog)
		if err != nil {
			queryError := apperrors.QueryError{
				Err:       err,
				Resource:  dog,
				Operation: apperrors.QueryUpdate,
			}
			return nil, &queryError
		}
		return createdDog, nil
	}
}

// NewAuthorizeFetchDogs authorizes fetching dogs
func NewAuthorizeFetchDogs() func(user models.User) (bool, error) {
	return func(user models.User) (bool, error) {
		return user.ID != "", nil
	}
}

// NewFetchDogs returns a service function for fetching dogs
func (f ServiceFactory) NewFetchDogs(
	authorize func(user models.User) (bool, error),
	fetch func(ctx context.Context) (*models.Dogs, error),
) func(ctx context.Context) (*models.Dogs, error) {
	return func(ctx context.Context) (*models.Dogs, error) {
		user, ok := appcontext.User(ctx)
		if !ok {
			contextError := apperrors.ContextError{
				Err:       errors.New("failed to get context"),
				Resource:  apperrors.ContextResourceUser,
				Operation: apperrors.ContextOperationGet,
			}
			return nil, &contextError
		}
		appcontext.LogRequestField(ctx, zap.String("user", user.ID))
		ok, err := authorize(user)
		if err != nil {
			appcontext.LogRequestError(ctx, "failed to authorize fetchDogs", err)
			return nil, err
		}
		if !ok {
			unauthorizedError := apperrors.UnauthorizedError{
				User:      user,
				Operation: apperrors.QueryFetch,
				Resource:  models.Dogs{},
				Err:       err,
			}
			return nil, &unauthorizedError
		}
		dogs, err := fetch(ctx)
		if err != nil {
			queryError := apperrors.QueryError{
				Err:       err,
				Resource:  models.Dogs{},
				Operation: apperrors.QueryFetch,
			}
			return nil, &queryError
		}
		return dogs, nil
	}
}
