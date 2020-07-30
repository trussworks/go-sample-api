package sources

import (
	"bin/bork/pkg/models"
	"context"
)

// DogLister fetches all the dogs in the system
type DogLister interface {
	ListDogs(context.Context) (*models.Dogs, error)
}

// DogListerFunc is an adapoter for any func with an appropriate signature
// to satisfy the DogLister interface
type DogListerFunc func(context.Context) (*models.Dogs, error)

// ListDogs satisfies the DogLister.ListDogs interface function
func (fn DogListerFunc) ListDogs(ctx context.Context) (*models.Dogs, error) {
	return fn(ctx)
}
