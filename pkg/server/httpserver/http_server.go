package httpserver

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"bin/bork/pkg/sources/postgres"
	"bin/bork/pkg/appconfig"
)

// Server holds dependencies for running the bork server extracted
// from appconfig
type Server struct {
	db        *sqlx.DB
	router    *mux.Router
	url       url.URL
	logger    *zap.Logger
	datetime  string
	version   string
	timestamp string
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// NewServer sets up the dependencies for a server
func NewServer(appConfig *appconfig.AppConfig) (*Server, error) {
	db, err := postgres.NewDB(appConfig)
	if err != nil {
		return nil, err
	}

	// Set the router
	r := mux.NewRouter()

	s := &Server{
		db:        db,
		router:    r,
		url:       NewURL(appConfig),
		logger:    appConfig.Logger,
		datetime:  appConfig.AppDatetime,
		version:   appConfig.AppVersion,
		timestamp: appConfig.AppTimestamp,
	}

	s.routes()

	return s, nil
}

// Serve runs the server
func Serve(appConfig *appconfig.AppConfig) {
	s, err := NewServer(appConfig)
	if err != nil {
		appConfig.Logger.Fatal(fmt.Sprintf("Failed to initialize server: %s", err))
		return
	}
	// start the server
	s.logger.Info("Serving application", zap.String("host", s.url.Host))
	err = http.ListenAndServe(s.url.Host, s)
	if err != nil {
		s.logger.Fatal(fmt.Sprintf("Failed to start server: %s", err))
	}
}
