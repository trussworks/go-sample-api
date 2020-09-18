package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"bin/bork/graph/generated"
	"bin/bork/graph/model"
	"bin/bork/pkg/appcontext"
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (r *mutationResolver) CreateDog(ctx context.Context, input model.DogInput) (*model.Dog, error) {
	user, ok := appcontext.User(ctx)
	if !ok {
		r.Logger.Error("Failed to get User from context")
		return nil, fmt.Errorf("Failed to get User from context")
	}
	dbInputDog, err := model.GqlDogInputToDbDog(input)
	if err != nil {
		return nil, err
	}
	ok, err = r.AuthorizeCreateDog(user, dbInputDog)
	if err != nil {
		r.Logger.Error("failed to authorize GQL CreateDog",
			zap.String("user", user.ID))
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("Unauthorized GQL CreateDog")
	}
	dbInputDog.OwnerID = user.ID
	dbDog, err := r.Store.CreateDog(ctx, dbInputDog)
	if err != nil {
		return nil, err
	}
	gqlDog, err := model.DbDogToGqlDog(dbDog)
	if err != nil {
		return nil, err
	}
	return gqlDog, nil
}

func (r *mutationResolver) UpdateDog(ctx context.Context, id string, input model.DogInput) (*model.Dog, error) {
	user, ok := appcontext.User(ctx)
	if !ok {
		r.Logger.Error("Failed to get User from context")
		return nil, fmt.Errorf("Failed to get User from context")
	}
	dbid := uuid.MustParse(id)
	dbDog, err := r.Store.FetchDog(ctx, dbid)
	if err != nil {
		return nil, err
	}
	ok, err = r.AuthorizeUpdateDog(user, dbDog)
	if err != nil {
		r.Logger.Error("failed to authorize GQL UpdateDog",
			zap.String("user", user.ID))
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("Unauthorized GQL UpdateDog")
	}
	dbInputDog, err := model.GqlDogInputToDbDog(input)
	if err != nil {
		return nil, err
	}
	dbDog.BirthDate = dbInputDog.BirthDate
	dbDog.Breed = dbInputDog.Breed
	dbDog.Name = dbInputDog.Name

	dbDog, err = r.Store.UpdateDog(ctx, dbDog)
	if err != nil {
		return nil, err
	}
	gqlDog, err := model.DbDogToGqlDog(dbDog)
	if err != nil {
		return nil, err
	}
	return gqlDog, nil
}

func (r *mutationResolver) Login(ctx context.Context, userID string) (*model.Owner, error) {
	// In the real world, validate credentials ...
	sessionCreator, ok := appcontext.SessionCreator(ctx)
	if !ok {
		return nil, fmt.Errorf("Cannot create session after login")
	}
	err := sessionCreator(userID)
	if err != nil {
		return nil, err
	}
	return &model.Owner{ID: userID}, nil
}

func (r *mutationResolver) Logout(ctx context.Context, userID string) (*model.Owner, error) {
	// must be logged in to log out
	user, ok := appcontext.User(ctx)
	if !ok {
		r.Logger.Error("Failed to get User from context")
		return nil, fmt.Errorf("Failed to get User from context")
	}
	r.Logger.Info(fmt.Sprintf("logging out user: %s", user))
	sessionCreator, ok := appcontext.SessionCreator(ctx)
	if !ok {
		return nil, fmt.Errorf("Cannot create session after login")
	}
	// fake invalidating session
	err := sessionCreator("")
	if err != nil {
		return nil, err
	}
	return &model.Owner{ID: user.ID, Email: user.Email}, nil
}

func (r *queryResolver) Dogs(ctx context.Context) ([]*model.Dog, error) {
	user, ok := appcontext.User(ctx)
	if !ok {
		r.Logger.Error("Failed to get User from context")
		return nil, fmt.Errorf("Failed to get User from context")
	}
	ok, err := r.AuthorizeFetchDogs(user)
	if err != nil {
		r.Logger.Error("failed to authorize GQL Dogs",
			zap.String("user", user.ID))
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("Unauthorized GQL Dogs")
	}
	dbDogs, err := r.Store.FetchDogs(ctx)
	if err != nil {
		return nil, err
	}
	allDogs := *dbDogs
	gqlDogs := make([]*model.Dog, len(allDogs))
	for i := range allDogs {
		dbDog := allDogs[i]
		gqlDogs[i], err = model.DbDogToGqlDog(&dbDog)
		if err != nil {
			return nil, err
		}
	}
	return gqlDogs, nil
}

func (r *queryResolver) Dog(ctx context.Context, dogID string) (*model.Dog, error) {
	user, ok := appcontext.User(ctx)
	if !ok {
		r.Logger.Error("Failed to get User from context")
		return nil, fmt.Errorf("Failed to get User from context")
	}
	uuid := uuid.MustParse(dogID)
	dbDog, err := r.Store.FetchDog(ctx, uuid)
	if err != nil {
		return nil, err
	}
	ok, err = r.AuthorizeFetchDog(user, dbDog)
	if err != nil {
		r.Logger.Error("failed to authorize GQL FetchDog",
			zap.String("user", user.ID))
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("Unauthorized GQL FetchDog")
	}
	gqlDog, err := model.DbDogToGqlDog(dbDog)
	if err != nil {
		return nil, err
	}
	return gqlDog, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
