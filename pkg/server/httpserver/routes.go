package httpserver

import (
	"github.com/google/uuid"

	"bin/bork/pkg/apis/v1/http/handlers"
	"bin/bork/pkg/models"
	"bin/bork/pkg/services"
)

func (s *Server) routes() {
	api := s.router.PathPrefix("/api/v1").Subrouter()

	handlerBase := handlers.NewHandlerBase(s.logger)

	serviceFactory := services.NewServiceFactory(s.logger)
	fakeAuthorize := func(user models.User, dog *models.Dog) (bool, error) {
		return true, nil
	}
	fakeFetchDog := func(id uuid.UUID) (*models.Dog, error) {
		return &models.Dog{}, nil
	}

	// endpoint for dog
	dogHandler := handlers.NewDogHandler(handlerBase, serviceFactory.NewFetchDog(fakeAuthorize, fakeFetchDog))
	api.Handle("/dog/{dog_id}", dogHandler.Handle())
}
