package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/facebookgo/clock"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"bin/bork/pkg/appcontext"
	"bin/bork/pkg/apperrors"
)

// NewHandlerBase is a constructor for HandlerBase
func NewHandlerBase(logger *zap.Logger) HandlerBase {
	return HandlerBase{
		logger: logger,
		clock:  clock.New(),
	}
}

// HandlerBase is for shared handler utilities
type HandlerBase struct {
	logger *zap.Logger
	clock  clock.Clock
}

type errorItem struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// errorResponse contains the structure of error for a http response
type errorResponse struct {
	Errors  []errorItem `json:"errors"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	TraceID uuid.UUID   `json:"traceID"`
}

func newErrorResponse(code int, message string, traceID uuid.UUID) errorResponse {
	return errorResponse{
		Errors:  []errorItem{},
		Code:    code,
		Message: message,
		TraceID: traceID,
	}
}

func (r *errorResponse) withValidations(validations apperrors.Validations) {
	for k, v := range validations {
		r.Errors = append(r.Errors, errorItem{
			Field:   k,
			Message: v,
		})
	}
}

// WriteErrorResponse writes a response for a given application error
func (b HandlerBase) WriteErrorResponse(ctx context.Context, w http.ResponseWriter, appErr error) {
	logger, ok := appcontext.Logger(ctx)
	if !ok {
		logger = b.logger
	}

	traceID, ok := appcontext.Trace(ctx)
	if !ok {
		traceID = uuid.New()
		logger.With(zap.String("traceID", traceID.String()))
	}

	// get code and response
	var code int
	var response errorResponse
	switch appErr := appErr.(type) {
	case *apperrors.UnauthorizedError:
		// 4XX errors are not logged as errors, but are for client
		logger.Info("Returning unauthorized response from handler", zap.Error(appErr))
		code = http.StatusUnauthorized
		response = newErrorResponse(
			http.StatusUnauthorized,
			"Unauthorized",
			traceID,
		)
	case *apperrors.QueryError:
		logger.Error("Returning server error response from handler", zap.Error(appErr))
		code = http.StatusInternalServerError
		response = newErrorResponse(
			http.StatusInternalServerError,
			"Something went wrong",
			traceID,
		)
	case *apperrors.ContextError:
		logger.Error("Returning server error response from handler", zap.Error(appErr))
		code = http.StatusInternalServerError
		response = newErrorResponse(
			http.StatusInternalServerError,
			"Something went wrong",
			traceID,
		)
	case *apperrors.ValidationError:
		logger.Info("Returning bad request error from handler", zap.Error(appErr))
		code = http.StatusBadRequest
		response = newErrorResponse(
			http.StatusBadRequest,
			"Bad request",
			traceID,
		)
		response.withValidations(appErr.Validations)
	case *apperrors.MethodNotAllowedError:
		logger.Info("Returning method not allowed error from handler", zap.Error(appErr))
		code = http.StatusMethodNotAllowed
		response = newErrorResponse(
			http.StatusMethodNotAllowed,
			"Method not allowed",
			traceID,
		)
	default:
		logger.Error("Returning server error response from handler", zap.Error(appErr))
		code = http.StatusInternalServerError
		response = newErrorResponse(
			http.StatusInternalServerError,
			"Something went wrong",
			traceID,
		)
	}

	// get error as response body
	responseBody, err := json.Marshal(response)
	if err != nil {
		logger.Error("Failed to marshal error response. Defaulting to generic.")
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// write a JSON response and fallback to generic message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(responseBody)
	if err != nil {
		logger.Error("Failed to write error response. Defaulting to generic.")
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}
