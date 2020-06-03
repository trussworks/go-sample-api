package httpserver

import (
	"github.com/google/uuid"

	"bin/bork/pkg/apis/v1/http/handlers"
	"bin/bork/pkg/models"
	httpserver "bin/bork/pkg/server/httpserver/middleware"
	"bin/bork/pkg/services"
)

func (s *Server) routes() {
	// health check goes directly on the main router to avoid auth
	healthCheckHandler := handlers.HealthCheckHandler{
		Config: s.Config,
	}
	s.router.HandleFunc("/api/v1/healthcheck", healthCheckHandler.Handle())

	api := s.router.PathPrefix("/api/v1").Subrouter()

	// set up base handler
	handlerBase := handlers.NewHandlerBase(s.logger)

	// use authorization on API
	api.Use(httpserver.NewFakeAuthorizeMiddleware(handlerBase))

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
