package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"bin/bork/pkg/appcontext"
	"bin/bork/pkg/apperrors"
	"bin/bork/pkg/models"
)

type fetchDog func(ctx context.Context, id uuid.UUID) (*models.Dog, error)

// NewDogHandler is a constructor for a DogHandler
func NewDogHandler(base HandlerBase, fetch fetchDog) DogHandler {
	return DogHandler{
		HandlerBase: base,
		fetchDog:    fetch,
	}
}

// DogHandler is the handler for CRUD operations on dog resources
type DogHandler struct {
	HandlerBase
	fetchDog fetchDog
}

// Handle handles a request for the system intake form
func (h DogHandler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger, ok := appcontext.Logger(r.Context())
		if !ok {
			h.logger.Error("Failed to get logger from context in system intake handler")
			logger = h.logger
		}

		switch r.Method {
		case "GET":
			idString := mux.Vars(r)["dog_id"]
			if idString == "" {
				validationErr := apperrors.NewValidationError(
					errors.New("GET dog params failed validation"),
					models.Dog{},
					idString,
				)
				validationErr.WithValidation("dogID", "required")
				h.WriteErrorResponse(
					r.Context(),
					w,
					&validationErr,
				)
				return
			}
			id, err := uuid.Parse(idString)
			if err != nil {
				validationErr := apperrors.NewValidationError(
					errors.New("GET dog params failed validation"),
					models.Dog{},
					idString,
				)
				validationErr.WithValidation("dogID", "must be UUID")
				h.WriteErrorResponse(
					r.Context(),
					w,
					&validationErr,
				)
				return
			}

			dog, err := h.fetchDog(r.Context(), id)
			if err != nil {
				h.WriteErrorResponse(r.Context(), w, err)
				return
			}

			responseBody, err := json.Marshal(dog)
			if err != nil {
				h.WriteErrorResponse(r.Context(), w, err)
				return
			}

			_, err = w.Write(responseBody)
			if err != nil {
				h.WriteErrorResponse(r.Context(), w, err)
				return
			}

			return

		default:
			logger.Info("Unsupported method requested")
			http.Error(w, "Method not allowed for system intake", http.StatusMethodNotAllowed)
			return
		}
	}
}
