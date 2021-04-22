package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"bin/bork/pkg/appcontext"
	"bin/bork/pkg/models"
)

func (s HandlerTestSuite) TestDogHandler_Handle() {
	dog := models.Dog{
		ID:        uuid.New(),
		Name:      "Chihua",
		Breed:     models.Chihuahua,
		BirthDate: s.base.clock.Now(),
	}
	fakeFetchDog := func(ctx context.Context, id uuid.UUID) (*models.Dog, error) {
		return &dog, nil
	}
	fakeCreateDog := func(ctx context.Context, dog *models.Dog) (*models.Dog, error) {
		return dog, nil
	}
	fakeUpdateDog := func(ctx context.Context, dog *models.Dog) (*models.Dog, error) {
		return dog, nil
	}

	requestContext := context.Background()
	requestContext = appcontext.WithUser(requestContext, models.User{ID: "McName"})
	s.Run("golden path GET returns 200", func() {
		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(
			requestContext,
			"GET",
			fmt.Sprintf("/dog/%s", dog.ID.String()),
			bytes.NewBufferString(""),
		)
		s.NoError(err)
		req = mux.SetURLVars(req, map[string]string{"dog_id": dog.ID.String()})

		DogHandler{
			s.base,
			fakeFetchDog,
			fakeCreateDog,
			fakeUpdateDog,
		}.Handle()(rr, req)

		s.Equal(http.StatusOK, rr.Code)
		responseDog := &models.Dog{}
		err = json.Unmarshal(rr.Body.Bytes(), responseDog)
		s.NoError(err)
		s.Equal(dog.ID, responseDog.ID)
	})

	s.Run("GET wth no ID returns 400", func() {
		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(
			requestContext,
			"GET",
			fmt.Sprintf("/dog/%s", dog.ID.String()),
			bytes.NewBufferString(""),
		)
		s.NoError(err)
		req = mux.SetURLVars(req, map[string]string{"dog_id": ""})

		DogHandler{
			s.base,
			fakeFetchDog,
			fakeCreateDog,
			fakeUpdateDog,
		}.Handle()(rr, req)

		s.Equal(http.StatusUnprocessableEntity, rr.Code)
	})

	s.Run("GET with bad ID returns 400", func() {
		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(
			requestContext,
			"GET",
			fmt.Sprintf("/dog/%s", dog.ID.String()),
			bytes.NewBufferString(""),
		)
		s.NoError(err)
		req = mux.SetURLVars(req, map[string]string{"dog_id": "badID"})

		DogHandler{
			s.base,
			fakeFetchDog,
			fakeCreateDog,
			fakeUpdateDog,
		}.Handle()(rr, req)

		s.Equal(http.StatusUnprocessableEntity, rr.Code)
	})

	s.Run("GET with fetch failing returns 500", func() {
		failFetchDog := func(ctx context.Context, id uuid.UUID) (*models.Dog, error) {
			return nil, errors.New("failed to fetch dog")
		}
		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(
			requestContext,
			"GET",
			fmt.Sprintf("/dog/%s", dog.ID.String()),
			bytes.NewBufferString(""),
		)
		s.NoError(err)
		req = mux.SetURLVars(req, map[string]string{"dog_id": dog.ID.String()})
		req = req.WithContext(appcontext.WithEmptyRequestLog(req.Context()))

		DogHandler{
			s.base,
			failFetchDog,
			fakeCreateDog,
			fakeUpdateDog,
		}.Handle()(rr, req)

		s.Equal(http.StatusInternalServerError, rr.Code)
	})

	s.Run("golden path POST returns 200", func() {
		rr := httptest.NewRecorder()
		dogBytes, err := json.Marshal(dog)
		s.NoError(err)
		req, err := http.NewRequestWithContext(
			requestContext,
			"POST",
			"/dog",
			bytes.NewBuffer(dogBytes),
		)
		s.NoError(err)

		DogHandler{
			s.base,
			fakeFetchDog,
			fakeCreateDog,
			fakeUpdateDog,
		}.Handle()(rr, req)

		s.Equal(http.StatusOK, rr.Code)
		responseDog := &models.Dog{}
		err = json.Unmarshal(rr.Body.Bytes(), responseDog)
		s.NoError(err)
		s.Equal(dog.ID, responseDog.ID)
	})

	s.Run("no body POST returns 400", func() {
		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(
			requestContext,
			"POST",
			"/dog",
			nil,
		)
		req = req.WithContext(appcontext.WithEmptyRequestLog(req.Context()))
		s.NoError(err)

		DogHandler{
			s.base,
			fakeFetchDog,
			fakeCreateDog,
			fakeUpdateDog,
		}.Handle()(rr, req)

		s.Equal(http.StatusBadRequest, rr.Code)
	})

	s.Run("bad body POST returns 400", func() {
		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(
			requestContext,
			"POST",
			"/dog",
			bytes.NewBufferString("{x: nil}"),
		)
		s.NoError(err)
		req = req.WithContext(appcontext.WithEmptyRequestLog(req.Context()))

		DogHandler{
			s.base,
			fakeFetchDog,
			fakeCreateDog,
			fakeUpdateDog,
		}.Handle()(rr, req)

		s.Equal(http.StatusBadRequest, rr.Code)
	})

	s.Run("when create fails POST returns 500", func() {
		failCreateDog := func(ctx context.Context, dog *models.Dog) (*models.Dog, error) {
			return nil, errors.New("failed to create dog")
		}
		rr := httptest.NewRecorder()
		dogBytes, err := json.Marshal(dog)
		s.NoError(err)
		req, err := http.NewRequestWithContext(
			requestContext,
			"POST",
			"/dog",
			bytes.NewBuffer(dogBytes),
		)
		s.NoError(err)
		req = req.WithContext(appcontext.WithEmptyRequestLog(req.Context()))

		DogHandler{
			s.base,
			fakeFetchDog,
			failCreateDog,
			fakeUpdateDog,
		}.Handle()(rr, req)

		s.Equal(http.StatusInternalServerError, rr.Code)
	})

	s.Run("golden path PUT returns 200", func() {
		rr := httptest.NewRecorder()
		dogBytes, err := json.Marshal(dog)
		s.NoError(err)
		req, err := http.NewRequestWithContext(
			requestContext,
			"PUT",
			"/dog",
			bytes.NewBuffer(dogBytes),
		)
		s.NoError(err)
		req = req.WithContext(appcontext.WithEmptyRequestLog(req.Context()))

		DogHandler{
			s.base,
			fakeFetchDog,
			fakeCreateDog,
			fakeUpdateDog,
		}.Handle()(rr, req)

		s.Equal(http.StatusOK, rr.Code)
		responseDog := &models.Dog{}
		err = json.Unmarshal(rr.Body.Bytes(), responseDog)
		s.NoError(err)
		s.Equal(dog.ID, responseDog.ID)
	})

	s.Run("no body PUT returns 400", func() {
		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(
			requestContext,
			"PUT",
			"/dog",
			nil,
		)
		s.NoError(err)
		req = req.WithContext(appcontext.WithEmptyRequestLog(req.Context()))

		DogHandler{
			s.base,
			fakeFetchDog,
			fakeCreateDog,
			fakeUpdateDog,
		}.Handle()(rr, req)

		s.Equal(http.StatusBadRequest, rr.Code)
	})

	s.Run("bad body PUT returns 400", func() {
		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(
			requestContext,
			"POST",
			"/dog",
			bytes.NewBufferString("{x: nil}"),
		)
		s.NoError(err)
		req = req.WithContext(appcontext.WithEmptyRequestLog(req.Context()))

		DogHandler{
			s.base,
			fakeFetchDog,
			fakeCreateDog,
			fakeUpdateDog,
		}.Handle()(rr, req)

		s.Equal(http.StatusBadRequest, rr.Code)
	})

	s.Run("when update fails PUT returns 500", func() {
		failUpdateDog := func(ctx context.Context, dog *models.Dog) (*models.Dog, error) {
			return nil, errors.New("failed to create dog")
		}
		rr := httptest.NewRecorder()
		dogBytes, err := json.Marshal(dog)
		s.NoError(err)
		req, err := http.NewRequestWithContext(
			requestContext,
			"PUT",
			"/dog",
			bytes.NewBuffer(dogBytes),
		)
		s.NoError(err)
		req = req.WithContext(appcontext.WithEmptyRequestLog(req.Context()))

		DogHandler{
			s.base,
			fakeFetchDog,
			fakeCreateDog,
			failUpdateDog,
		}.Handle()(rr, req)

		s.Equal(http.StatusInternalServerError, rr.Code)
	})

	s.Run("unsupported method returns 405", func() {
		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(
			requestContext,
			"OPTIONS",
			fmt.Sprintf("/dog/%s", dog.ID.String()),
			bytes.NewBufferString(""),
		)
		s.NoError(err)
		req = req.WithContext(appcontext.WithEmptyRequestLog(req.Context()))
		req = mux.SetURLVars(req, map[string]string{"dog_id": dog.ID.String()})

		DogHandler{
			s.base,
			fakeFetchDog,
			fakeCreateDog,
			fakeUpdateDog,
		}.Handle()(rr, req)

		s.Equal(http.StatusMethodNotAllowed, rr.Code)
	})
}
