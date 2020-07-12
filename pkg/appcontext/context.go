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
	requestLogKey
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
func WithTrace(ctx context.Context) (context.Context, uuid.UUID) {
	traceID := uuid.New()
	return context.WithValue(ctx, traceKey, traceID), traceID
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

// requestLog keeps track of all the info to be logged in the request log line
type requestLog struct {
	fields       []zap.Field
	errorMessage string
	error        bool
}

// LogRequestField adds a zap.Field to the line logged at the end of the request
func LogRequestField(ctx context.Context, field zap.Field) {
	requestLog, ok := ctx.Value(requestLogKey).(*requestLog)
	if !ok {
		panic("Configuration Error, make sure you call WithEmptyRequestLog before LogRequestField")
	}

	// TODO, figure out how to make this threadsafe
	requestLog.fields = append(requestLog.fields, field)
}

// LogRequestError adds a message to the request log line and also sets it to log at the Error level
func LogRequestError(ctx context.Context, message string, err error) {
	requestLog, ok := ctx.Value(requestLogKey).(*requestLog)
	if !ok {
		panic("Configuration Error, make sure you call WithEmptyRequestLog before LogRequestField")
	}

	requestLog.fields = append(requestLog.fields, zap.Error(err))
	requestLog.errorMessage = message
	requestLog.error = true
}

// WithEmptyRequestLog returns a context with the request User
func WithEmptyRequestLog(ctx context.Context) context.Context {
	return context.WithValue(ctx, requestLogKey, &requestLog{})
}

// RequestLogFields returns all the fields added during the request
func RequestLogFields(ctx context.Context) ([]zap.Field, bool) {
	requestLog, ok := ctx.Value(requestLogKey).(*requestLog)

	return requestLog.fields, ok
}
