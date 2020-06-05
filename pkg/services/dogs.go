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
	fetch func(id uuid.UUID) (*models.Dog, error),
) func(ctx context.Context, id uuid.UUID) (*models.Dog, error) {
	return func(ctx context.Context, id uuid.UUID) (*models.Dog, error) {
		logger, ok := appcontext.Logger(ctx)
		if !ok {
			logger = f.logger
			logger.Error("failed to get logger from context in FetchDog service")
		}
		dog, err := fetch(id)
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
		ok, err = authorize(user, dog)
		if err != nil {
			logger.Error("failed to authorize fetchDog", zap.String("user", user.ID))
			return nil, err
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
		if user.ID != "" {
			return true, nil
		}
		return false, nil
	}
}

// NewCreateDog returns a service function for creating a dog
func (f ServiceFactory) NewCreateDog(
	authorize func(user models.User, dog *models.Dog) (bool, error),
	create func(dog *models.Dog) (*models.Dog, error),
) func(ctx context.Context, dog *models.Dog) (*models.Dog, error) {
	return func(ctx context.Context, dog *models.Dog) (*models.Dog, error) {
		logger, ok := appcontext.Logger(ctx)
		if !ok {
			logger = f.logger
			logger.Error("failed to get logger from context in CreateDog service")
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
		ok, err := authorize(user, dog)
		if err != nil {
			logger.Error("failed to authorize createDog", zap.String("user", user.ID))
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
		createdDog, err := create(dog)
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
	update func(dog *models.Dog) (*models.Dog, error),
	fetch func(id uuid.UUID) (*models.Dog, error),
) func(ctx context.Context, dog *models.Dog) (*models.Dog, error) {
	return func(ctx context.Context, dog *models.Dog) (*models.Dog, error) {
		logger, ok := appcontext.Logger(ctx)
		if !ok {
			logger = f.logger
			logger.Error("failed to get logger from context in UpdateDog service")
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
		existingDog, err := fetch(dog.ID)
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
			logger.Error("failed to authorize updateDog", zap.String("user", user.ID))
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
		createdDog, err := update(dog)
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
