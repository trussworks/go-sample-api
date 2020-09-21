package appcontext

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"bin/bork/pkg/models"
)

type contextKey int

const (
	loggerKey contextKey = iota
	traceKey
	userKey
	sessionCreatorKey
)

// WithLogger returns a context with the given logger
func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// Logger returns the context's logger
func Logger(ctx context.Context) (*zap.Logger, bool) {
	logger, ok := ctx.Value(loggerKey).(*zap.Logger)
	return logger, ok
}

// WithTrace returns a context with request trace
func WithTrace(ctx context.Context) context.Context {
	traceID := uuid.New()
	return context.WithValue(ctx, traceKey, traceID)
}

// Trace returns the context's trace UUID
func Trace(ctx context.Context) (uuid.UUID, bool) {
	traceID, ok := ctx.Value(traceKey).(uuid.UUID)
	return traceID, ok
}

// WithUser returns a context with the request User
func WithUser(ctx context.Context, user models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// User returns the context's User
func User(ctx context.Context) (models.User, bool) {
	user, ok := ctx.Value(userKey).(models.User)
	return user, ok
}

// SessionCreatorFunc creates a session for the client
type SessionCreatorFunc func(session string) error

// WithSessionCreator returns a context with the session creator function
func WithSessionCreator(ctx context.Context, sessionCreator SessionCreatorFunc) context.Context {
	return context.WithValue(ctx, sessionCreatorKey, sessionCreator)
}

func SessionCreator(ctx context.Context) (SessionCreatorFunc, bool) {
	sessionCreator, ok := ctx.Value(sessionCreatorKey).(SessionCreatorFunc)
	return sessionCreator, ok
}
