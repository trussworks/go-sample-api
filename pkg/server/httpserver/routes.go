package httpserver

import (
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"bin/bork/graph"
	"bin/bork/graph/generated"
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
	handlerBase := handlers.NewHandlerBase(s.logger, s.clock)

	// health check goes directly on the main router to avoid auth
	healthCheckHandler := handlers.NewHealthCheckHandler(
		handlerBase,
		handlers.HealthDatetime(s.datetime),
		handlers.HealthVersion(s.version),
		handlers.HealthTimestamp(s.timestamp),
	)
	s.router.HandleFunc("/api/v1/healthcheck", healthCheckHandler.Handle())

	// add a request based logger
	api.Use(NewLoggerMiddleware(s.logger))

	// use authorization on API
	api.Use(NewFakeAuthorizeMiddleware(handlerBase))

	// set up service factory
	serviceFactory := services.NewServiceFactory(s.logger, s.clock)

	// create store
	store := postgres.NewStoreWithDB(s.db)

	fetchDogService := serviceFactory.NewFetchDog(
		services.NewAuthorizeFetchDog(),
		store.FetchDog,
	)
	createDogService := serviceFactory.NewCreateDog(
		services.NewAuthorizeCreateDog(),
		store.CreateDog,
	)
	updateDogService := serviceFactory.NewUpdateDog(
		services.NewAuthorizeUpdateDog(),
		store.UpdateDog,
		store.FetchDog,
	)

	dogHandler := handlers.NewDogHandler(
		handlerBase,
		fetchDogService,
		createDogService,
		updateDogService,
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
	fetchDogsService := serviceFactory.NewFetchDogs(
		services.NewAuthorizeFetchDogs(),
		cache.NewStore(cacheConfig).FetchDogs,
	)

	dogsHandler := handlers.NewDogsHandler(
		handlerBase,
		fetchDogsService,
	)
	api.Handle("/dogs", dogsHandler.Handle())

	gRoute := s.router.PathPrefix("/graphql").Subrouter()
	// add a request based logger
	gRoute.Use(NewLoggerMiddleware(s.logger))
	// Cookies for GraphQL
	gRoute.Use(NewSimpleSessionMiddleware(handlerBase, "/graphql"))

	gRoute.Handle("/playground", playground.Handler("GraphQL playground",
		"/graphql/query"))

	graphResolver := &graph.Resolver{
		Clock:       s.clock,
		Logger:      s.logger,
		FetchDbDog:  fetchDogService,
		CreateDbDog: createDogService,
		UpdateDbDog: updateDogService,
		FetchDbDogs: fetchDogsService,
	}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(
		generated.Config{Resolvers: graphResolver}))
	gRoute.Handle("/query", srv)

	s.router.PathPrefix("/").Handler(handlers.NewCatchAllHandler(handlerBase).Handle())
}
