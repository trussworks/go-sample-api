package httpserver

import (
	"github.com/google/uuid"

	"bin/bork/pkg/apis/v1/http/handlers"
	"bin/bork/pkg/models"
	"bin/bork/pkg/services"
)

func (s *Server) routes() {
	// trace all requests with an ID
	s.router.Use(NewTraceMiddleware())

	api := s.router.PathPrefix("/api/v1").Subrouter()

	// set up base handler
	handlerBase := handlers.NewHandlerBase(s.logger)

	// health check goes directly on the main router to avoid auth
	healthCheckHandler := handlers.NewHealthCheckHandler(
		handlerBase,
		s.Config,
	)
	s.router.HandleFunc("/api/v1/healthcheck", healthCheckHandler.Handle())

	// add a request based logger
	api.Use(NewLoggerMiddleware(s.logger))

	// use authorization on API
	api.Use(NewFakeAuthorizeMiddleware(handlerBase))

	// set up service factory
	serviceFactory := services.NewServiceFactory(s.logger)

	// endpoint for dog
	fakeAuthorize := func(user models.User, dog *models.Dog) (bool, error) {
		return true, nil
	}
	fakeFetchDog := func(id uuid.UUID) (*models.Dog, error) {
		return &models.Dog{}, nil
	}

	dogHandler := handlers.NewDogHandler(handlerBase, serviceFactory.NewFetchDog(fakeAuthorize, fakeFetchDog))
	api.Handle("/dog/{dog_id}", dogHandler.Handle())
}
