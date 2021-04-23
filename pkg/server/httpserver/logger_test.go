package httpserver

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"bin/bork/pkg/appcontext"
)

func (s ServerTestSuite) TestLoggerMiddleware() {
	s.Run("do a single log field", func() {

		req := httptest.NewRequest("GET", "/dogs/", nil)
		rr := httptest.NewRecorder()
		traceMiddleware := NewTraceMiddleware()

		// let's get a logger we can inspect
		logcore, recorded := observer.New(zapcore.InfoLevel)
		testLogger := zap.New(logcore)

		loggerMiddleware := NewLoggerMiddleware(testLogger)

		// The handler will log a field and write a header
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// let's add some log fields
			appcontext.LogRequestField(r.Context(), zap.String("onekey", "onevalue"))

			w.WriteHeader(404)
		})

		traceMiddleware(loggerMiddleware(testHandler)).ServeHTTP(rr, req)

		// Check that we logged what we expected
		s.Equal(1, recorded.Len())

		line := recorded.All()[0]
		s.Equal("Request Complete", line.Message)

		s.Contains(line.Context, zap.String("host", "example.com"))
		s.Contains(line.Context, zap.String("onekey", "onevalue"))
		s.Contains(line.Context, zap.Int("http_status", 404))

	})

	s.Run("log an error", func() {

		req := httptest.NewRequest("GET", "/dogs/", nil)
		rr := httptest.NewRecorder()
		traceMiddleware := NewTraceMiddleware()

		// let's get a logger we can inspect
		logcore, recorded := observer.New(zapcore.InfoLevel)
		testLogger := zap.New(logcore)

		loggerMiddleware := NewLoggerMiddleware(testLogger)

		// The handler will log a field and write a header
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			err := errors.New("this test is in error")

			// let's add some log fields
			appcontext.LogRequestError(r.Context(), "the test errored just like we planned", err)

			w.WriteHeader(500)
		})

		traceMiddleware(loggerMiddleware(testHandler)).ServeHTTP(rr, req)

		// Check that we logged what we expected
		s.Equal(1, recorded.Len())

		line := recorded.All()[0]
		s.Equal("Request Complete", line.Message)
		s.Equal(zap.ErrorLevel, line.Level)

		s.Contains(line.Context, zap.String("host", "example.com"))
		s.Contains(line.Context, zap.Int("http_status", 500))
		s.Contains(line.Context, zap.String("error_message", "the test errored just like we planned"))

	})

}
