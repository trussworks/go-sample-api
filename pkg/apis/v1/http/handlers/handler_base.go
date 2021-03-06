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

func (r *errorResponse) withMap(errMap map[string]string) {
	for k, v := range errMap {
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
		code = http.StatusUnauthorized
		response = newErrorResponse(
			code,
			"Unauthorized",
			traceID,
		)
	case *apperrors.QueryError:
		switch appErr.Unwrap().(type) {
		case *apperrors.ResourceNotFoundError:
			code = http.StatusNotFound
			response = newErrorResponse(
				code,
				"Resource not found",
				traceID,
			)
		default:
			code = http.StatusInternalServerError
			appcontext.LogRequestError(ctx, "DB Query Error", appErr)
			response = newErrorResponse(
				code,
				"Something went wrong",
				traceID,
			)
		}
	case *apperrors.ContextError:
		code = http.StatusInternalServerError
		appcontext.LogRequestError(ctx, "Context Error", appErr)
		response = newErrorResponse(
			code,
			"Something went wrong",
			traceID,
		)
	case *apperrors.ValidationError:
		code = http.StatusUnprocessableEntity
		response = newErrorResponse(
			code,
			"Entity unprocessable",
			traceID,
		)
		response.withMap(appErr.Validations.Map())
	case *apperrors.MethodNotAllowedError:
		code = http.StatusMethodNotAllowed
		response = newErrorResponse(
			code,
			"Method not allowed",
			traceID,
		)
	case *apperrors.UnknownRouteError:
		logger.Info("Returning status not found error from handler", zap.Error(appErr))
		code = http.StatusNotFound
		response = newErrorResponse(
			code,
			"Not found",
			traceID,
		)
	case *apperrors.BadRequestError:
		code = http.StatusBadRequest
		appcontext.LogRequestField(ctx, zap.NamedError("bad_request_err", appErr))
		response = newErrorResponse(
			code,
			"Bad request",
			traceID,
		)
	default:
		code = http.StatusInternalServerError
		appcontext.LogRequestError(ctx, "Unexpected Error", appErr)
		response = newErrorResponse(
			code,
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
