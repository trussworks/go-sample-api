package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"bin/bork/pkg/appcontext"
)

//type fetchDog func(ctx context.Context, id uuid.UUID) (*models.Dog, error)

// NewDogHandler is a constructor for a DogHandler
func NewDogHandler(base HandlerBase) DogHandler {
	return DogHandler{
		base,
	}
}

// DogHandler is the handler for CRUD operations on dog resources
type DogHandler struct {
	HandlerBase
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
			id := mux.Vars(r)["dog_id"]
			if id == "" {
				http.Error(w, "Dog ID required", http.StatusBadRequest)
				return
			}
			_, err := uuid.Parse(id)
			if err != nil {
				logger.Error("Failed to parse dog id to uuid")
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			// Fetch Dog here
			// dog, err := h.FetchDogByID(uuid)
			//if err != nil {
			//	http.Error(w, "Failed to GET dog", http.StatusInternalServerError)
			//	return
			//}

			//responseBody, err := json.Marshal(intake)
			//if err != nil {
			//	http.Error(w, err.Error(), http.StatusInternalServerError)
			//	return
			//}

			//_, err = w.Write(responseBody)
			//if err != nil {
			//	http.Error(w, "Failed to get dog by id", http.StatusInternalServerError)
			//	return
			//}

			return

		default:
			logger.Info("Unsupported method requested")
			http.Error(w, "Method not allowed for system intake", http.StatusMethodNotAllowed)
			return
		}
	}
}
