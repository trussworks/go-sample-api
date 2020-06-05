package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"bin/bork/pkg/apperrors"
	"bin/bork/pkg/models"
)

type fetchDogs func(ctx context.Context) (*models.Dogs, error)

// NewDogsHandler is a constructor for a DogHandler
func NewDogsHandler(base HandlerBase, fetch fetchDogs) DogsHandler {
	return DogsHandler{
		HandlerBase: base,
		fetchDogs:   fetch,
	}
}

// DogsHandler is the handler for API operations on dog lists
type DogsHandler struct {
	HandlerBase
	fetchDogs fetchDogs
}

// Handle handles a request for a dog
func (h DogsHandler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			dogs, err := h.fetchDogs(r.Context())
			if err != nil {
				h.WriteErrorResponse(r.Context(), w, err)
				return
			}

			responseBody, err := json.Marshal(dogs)
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
		default:
			h.WriteErrorResponse(r.Context(), w, &apperrors.MethodNotAllowedError{Method: r.Method})
			return
		}
	}
}
