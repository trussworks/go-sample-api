package handlers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"bin/bork/pkg/appcontext"
	"bin/bork/pkg/models"
)

func (s HandlerTestSuite) TestDogHandler_Handle() {
	requestContext := context.Background()
	requestContext = appcontext.WithUser(requestContext, models.User{ID: "McName"})
	id, err := uuid.NewUUID()
	s.NoError(err)
	s.Run("golden path GET passes", func() {
		rr := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(requestContext, "GET", fmt.Sprintf("/dog/%s", id.String()), bytes.NewBufferString(""))
		s.NoError(err)
		req = mux.SetURLVars(req, map[string]string{"dog_id": id.String()})

		DogHandler{
			s.base,
		}.Handle()(rr, req)

		s.Equal(http.StatusOK, rr.Code)
	})
}
