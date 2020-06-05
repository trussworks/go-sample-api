package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"

	"bin/bork/pkg/appcontext"
	"bin/bork/pkg/models"
)

func (s HandlerTestSuite) TestDogsHandler_Handle() {
	dog := models.Dog{
		ID:        uuid.New(),
		Name:      "Chihua",
		Breed:     models.Chihuahua,
		BirthDate: s.base.clock.Now(),
	}
	fakeFetchDogs := func(ctx context.Context) (*models.Dogs, error) {
		return &models.Dogs{dog}, nil
	}

	requestContext := context.Background()
	requestContext = appcontext.WithUser(requestContext, models.User{ID: "McName"})
	s.Run("golden path GET returns 200", func() {
		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(
			requestContext,
			"GET",
			"/dogs",
			bytes.NewBufferString(""),
		)
		s.NoError(err)

		DogsHandler{
			s.base,
			fakeFetchDogs,
		}.Handle()(rr, req)

		s.Equal(http.StatusOK, rr.Code)
		responseDogs := &models.Dogs{}
		err = json.Unmarshal(rr.Body.Bytes(), responseDogs)
		s.NoError(err)
		s.Len(*responseDogs, 1)
		s.Equal(dog.ID, (*responseDogs)[0].ID)
	})

	s.Run("GET with fetch failing returns 500", func() {
		failFetchDogs := func(ctx context.Context) (*models.Dogs, error) {
			return nil, errors.New("failed to fetch dog")
		}
		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(
			requestContext,
			"GET",
			"/dogs",
			bytes.NewBufferString(""),
		)
		s.NoError(err)

		DogsHandler{
			s.base,
			failFetchDogs,
		}.Handle()(rr, req)

		s.Equal(http.StatusInternalServerError, rr.Code)
	})

	s.Run("unsupported method returns 405", func() {
		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(
			requestContext,
			"OPTIONS",
			"/dogs",
			bytes.NewBufferString(""),
		)
		s.NoError(err)

		DogsHandler{
			s.base,
			fakeFetchDogs,
		}.Handle()(rr, req)

		s.Equal(http.StatusMethodNotAllowed, rr.Code)
	})
}
