package httpserver

import (
	"net/http"

	"go.uber.org/zap"

	"bin/bork/pkg/appcontext"
)

const traceField string = "traceID"

// In order to log the status we need to wrap the ResponseWriter
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func newStatusRecorder(w http.ResponseWriter) *statusRecorder {
	return &statusRecorder{w, 200}
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func loggerMiddleware(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		traceID, ok := appcontext.Trace(ctx)
		if ok {
			logger = logger.With(zap.String(traceField, traceID.String()))
		} else {
			logger.Error("Failed to get trace ID from context")
		}
		ctx = appcontext.WithEmptyRequestLog(ctx)

		fields := []zap.Field{
			zap.String("accepted-language", r.Header.Get("accepted-language")),
			zap.Int64("content-length", r.ContentLength),
			zap.String("host", r.Host),
			zap.String("method", r.Method),
			zap.String("protocol-version", r.Proto),
			zap.String("referer", r.Header.Get("referer")),
			zap.String("source", r.RemoteAddr),
			zap.String("url", r.URL.String()),
			zap.String("user-agent", r.UserAgent()),
		}

		rec := newStatusRecorder(w)
		next.ServeHTTP(rec, r.WithContext(ctx))

		// get a couple more default fields
		fields = append(fields, zap.Int("http_status", rec.status))

		requestFields, ok := appcontext.RequestLogFields(ctx)
		if !ok {
			logger.Error("Fields not configured for this request")
		}

		allfields := append(fields, requestFields...)

		didError, errorMessage, ok := appcontext.RequestErrorInfo(ctx)
		if !ok {
			logger.Error("Error info not found on this request")
		}

		if didError {
			allfields = append(allfields, zap.String("error_message", errorMessage))
			logger.Error("Request Complete", allfields...)
		} else {
			logger.Info("Request Complete", allfields...)
		}
	})
}

// NewLoggerMiddleware returns a handler with a request based logger
func NewLoggerMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return loggerMiddleware(logger, next)
	}
}
