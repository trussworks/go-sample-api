package httpserver

import (
	"time"

	"go.uber.org/zap"

	"bin/bork/pkg/apis/v1/http/handlers"
	"bin/bork/pkg/services"
	"bin/bork/pkg/sources/cache"
	"bin/bork/pkg/sources/memory"
	"bin/bork/pkg/sources/postgres"
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

	//	create store
	store, err := postgres.NewStore(
		s.NewDBConfig(),
	)
	if err != nil {
		s.logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	dogHandler := handlers.NewDogHandler(
		handlerBase,
		serviceFactory.NewFetchDog(
			services.NewAuthorizeFetchDog(),
			store.FetchDog,
		),
		serviceFactory.NewCreateDog(
			services.NewAuthorizeCreateDog(),
			store.CreateDog,
		),
		serviceFactory.NewUpdateDog(
			services.NewAuthorizeUpdateDog(),
			store.UpdateDog,
			store.FetchDog,
		),
	)
	api.Handle("/dog/{dog_id}", dogHandler.Handle())
	api.Handle("/dog", dogHandler.Handle())

	// fabricated example of a cache store
	// to show composable store patterns
	cacheConfig := cache.StoreConfig{
		TTL:           time.Minute,
		DogCacheStore: memory.NewStore(),
		DogReadStore:  store,
	}
	dogsHandler := handlers.NewDogsHandler(
		handlerBase,
		serviceFactory.NewFetchDogs(
			services.NewAuthorizeFetchDogs(),
			cache.NewStore(cacheConfig).FetchDogs,
		),
	)
	api.Handle("/dogs", dogsHandler.Handle())

	s.router.PathPrefix("/").Handler(handlers.NewCatchAllHandler(handlerBase).Handle())
}
