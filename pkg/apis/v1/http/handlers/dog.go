package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"bin/bork/pkg/apperrors"
	"bin/bork/pkg/models"
)

type fetchDog func(ctx context.Context, id uuid.UUID) (*models.Dog, error)
type createDog func(ctx context.Context, dog *models.Dog) (*models.Dog, error)

// NewDogHandler is a constructor for a DogHandler
func NewDogHandler(base HandlerBase, fetch fetchDog, create createDog) DogHandler {
	return DogHandler{
		HandlerBase: base,
		fetchDog:    fetch,
		createDog:   create,
	}
}

// DogHandler is the handler for CRUD operations on dog resources
type DogHandler struct {
	HandlerBase
	fetchDog  fetchDog
	createDog createDog
}

// Handle handles a request for a dog
func (h DogHandler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

			w.Header().Set("Content-Type", "application/json")
			_, err = w.Write(responseBody)
			if err != nil {
				h.WriteErrorResponse(r.Context(), w, err)
				return
			}

			return
		case "POST":
			if r.Body == nil {
				h.WriteErrorResponse(
					r.Context(),
					w,
					&apperrors.BadRequestError{
						Err: errors.New("empty request not allowed"),
					},
				)
				return
			}
			defer r.Body.Close()
			decoder := json.NewDecoder(r.Body)
			dog := models.Dog{}
			err := decoder.Decode(&dog)

			if err != nil {
				h.WriteErrorResponse(
					r.Context(),
					w,
					&apperrors.BadRequestError{
						Err: errors.New("unable to decode dog"),
					},
				)
				return
			}

			responseDog, err := h.createDog(r.Context(), &dog)
			if err != nil {
				h.WriteErrorResponse(r.Context(), w, err)
				return
			}

			responseBody, err := json.Marshal(responseDog)
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
			h.WriteErrorResponse(r.Context(), w, &apperrors.MethodNotAllowedError{Method: r.Method})
			return
		}
	}
}
