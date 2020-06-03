package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/facebookgo/clock"
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
}

// WriteErrorResponse writes a response for a given application error
func (b HandlerBase) WriteErrorResponse(ctx context.Context, w http.ResponseWriter, appErr error) {
	logger, ok := appcontext.Logger(ctx)
	if !ok {
		logger = b.logger
	}

	// get code and response
	var code int
	var response errorResponse
	switch appErr.(type) {
	case *apperrors.UnauthorizedError:
		// 4XX errors are not logged as errors, but are for client
		logger.Info("Returning unauthorized response from handler", zap.Error(appErr))
		code = http.StatusUnauthorized
		response = errorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
	case *apperrors.QueryError:
		logger.Error("Returning server error response from handler", zap.Error(appErr))
		code = http.StatusInternalServerError
		response = errorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong",
		}
	case *apperrors.ContextError:
		logger.Error("Returning unknown error response from handler", zap.Error(appErr))
		code = http.StatusInternalServerError
		response = errorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong",
		}
	default:
		logger.Error("Returning server error response from handler", zap.Error(appErr))
		code = http.StatusInternalServerError
		response = errorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong",
		}
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
