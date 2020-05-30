package httpserver

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"bin/bork/pkg/appconfig"
)

// Server holds dependencies for running the EASi server
type Server struct {
	router      *mux.Router
	Config      *viper.Viper
	logger      *zap.Logger
	environment appconfig.Environment
	url         url.URL
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// NewServer sets up the dependencies for a server
func NewServer(config *viper.Viper) *Server {
	// Set up logger first so we can use it
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Failed to initial logger.")
	}

	// Set environment from config
	environment, err := appconfig.NewEnvironment(config.GetString(appconfig.EnvironmentKey))
	if err != nil {
		zapLogger.Fatal("Unable to set environment", zap.Error(err))
	}

	// Set the router
	r := mux.NewRouter()

	// TODO: set up routes

	s := &Server{
		router:      r,
		Config:      config,
		logger:      zapLogger,
		environment: environment,
	}

	s.url = s.NewURL()

	return s
}

// Serve runs the server
func Serve(config *viper.Viper) {
	s := NewServer(config)
	// start the server
	s.logger.Info("Serving application", zap.String("host", s.url.Host))
	err := http.ListenAndServe(s.url.Host, s)
	if err != nil {
		s.logger.Fatal("Failed to start server")
	}
}
