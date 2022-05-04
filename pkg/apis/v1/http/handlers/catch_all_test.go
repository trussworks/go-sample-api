package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"bin/bork/pkg/appcontext"
)

func (s HandlerTestSuite) TestCatchAllHandler() {
	s.Run("catch all handler always returns 404", func() {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/notAURL", bytes.NewBufferString(""))
		s.NoError(err)
		req = req.WithContext(appcontext.WithEmptyRequestLog(req.Context()))
		CatchAllHandler{
			s.base,
		}.Handle()(rr, req)
		s.Equal(http.StatusNotFound, rr.Code)
	})
}
