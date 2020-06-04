package httpserver

import (
	"net/http"

	"bin/bork/pkg/appcontext"
)

func traceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		next.ServeHTTP(w, r.WithContext(appcontext.WithTrace(ctx)))
	})
}

// NewTraceMiddleware returns a handler with a trace ID in context
func NewTraceMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return traceMiddleware(next)
	}
}
