package httpserver

import "bin/bork/pkg/apis/v1/http/handlers"

func (s *Server) routes() {
	api := s.router.PathPrefix("/api/v1").Subrouter()

	handlerBase := handlers.NewHandlerBase(s.logger)

	// endpoint for dog
	dogHandler := handlers.NewDogHandler(handlerBase)
	api.Handle("/dog/{dog_id}", dogHandler.Handle())
}
