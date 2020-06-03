package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	"bin/bork/pkg/apperrors"
)

type failWriter struct {
	realWriter *httptest.ResponseRecorder
	failCount  int
}

func (w *failWriter) Write(b []byte) (int, error) {
	if w.failCount == 0 {
		return w.realWriter.Write(b)
	}
	w.failCount--
	return 0, errors.New("writer fails")
}
func (w *failWriter) WriteHeader(statusCode int) {
	w.realWriter.WriteHeader(statusCode)
}

func (w *failWriter) Header() http.Header {
	return w.realWriter.Header()
}

func (s HandlerTestSuite) TestWriteErrorResponse() {
	ctx := context.Background()

	var responseTests = []struct {
		appErr      error
		code        int
		errResponse errorResponse
	}{
		{
			&apperrors.UnauthorizedError{},
			http.StatusUnauthorized,
			errorResponse{
				Errors:  []errorItem{},
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
			},
		},
		{
			&apperrors.QueryError{},
			http.StatusInternalServerError,
			errorResponse{
				Errors:  []errorItem{},
				Code:    http.StatusInternalServerError,
				Message: "Something went wrong",
			},
		},
		{
			&apperrors.ContextError{},
			http.StatusInternalServerError,
			errorResponse{
				Errors:  []errorItem{},
				Code:    http.StatusInternalServerError,
				Message: "Something went wrong",
			},
		},
		{
			&apperrors.ValidationError{
				Validations: map[string]string{"key": "required"},
			},
			http.StatusBadRequest,
			errorResponse{
				Errors:  []errorItem{{Field: "key", Message: "required"}},
				Code:    http.StatusBadRequest,
				Message: "Bad request",
			},
		},
		{
			errors.New("unknown error"),
			http.StatusInternalServerError,
			errorResponse{
				Errors:  []errorItem{},
				Code:    http.StatusInternalServerError,
				Message: "Something went wrong",
			},
		},
	}
	for _, t := range responseTests {
		s.Run(fmt.Sprintf("%T returns %d code", t.appErr, t.code), func() {
			writer := httptest.NewRecorder()

			s.base.WriteErrorResponse(ctx, writer, t.appErr)

			s.Equal(t.code, writer.Code)
			s.Equal("application/json", writer.Header().Get("Content-Type"))
			errResponse := &errorResponse{}
			err := json.Unmarshal(writer.Body.Bytes(), errResponse)
			s.NoError(err)
			s.Equal(t.errResponse, *errResponse)
		})
	}

	s.Run("failing to write json return plain text response", func() {
		writer := failWriter{
			realWriter: httptest.NewRecorder(),
			failCount:  1,
		}
		s.base.WriteErrorResponse(ctx, &writer, errors.New("some error"))

		s.Equal("text/plain; charset=utf-8", writer.Header().Get("Content-Type"))
		s.Equal(http.StatusInternalServerError, writer.realWriter.Code)
		s.Equal("Something went wrong\n", writer.realWriter.Body.String())
	})
}
