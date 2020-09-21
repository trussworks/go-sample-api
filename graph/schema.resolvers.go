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
)

func (r *mutationResolver) CreateDog(ctx context.Context, input model.DogInput) (*model.Dog, error) {
	dbInputDog, err := model.GqlDogInputToDbDog(input)
	if err != nil {
		return nil, err
	}
	dbDog, err := r.CreateDbDog(ctx, dbInputDog)
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
	dbInputDog, err := model.GqlDogInputToDbDog(input)
	if err != nil {
		return nil, err
	}
	dbInputDog.ID = uuid.MustParse(id)
	dbDog, err := r.UpdateDbDog(ctx, dbInputDog)
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
	dbDogs, err := r.FetchDbDogs(ctx)
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
	dbid := uuid.MustParse(dogID)
	dbDog, err := r.FetchDbDog(ctx, dbid)
	if err != nil {
		return nil, err
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
