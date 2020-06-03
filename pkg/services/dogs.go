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