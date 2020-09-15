package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"bin/bork/graph/generated"
	"bin/bork/graph/model"
	"bin/bork/pkg/appcontext"
)

func (r *mutationResolver) CreateDog(ctx context.Context, input model.NewDog) (*model.Dog, error) {
	dbInputDog, err := model.GqlNewDogToDbDog(input)
	if err != nil {
		return nil, err
	}
	user, ok := appcontext.User(ctx)
	if !ok {
		r.Logger.Error("failed to get context")
		return nil, fmt.Errorf("failed to get context")
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

func (r *queryResolver) Dogs(ctx context.Context) ([]*model.Dog, error) {
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

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
