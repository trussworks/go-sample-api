package handlers

import (
	"context"
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
		appErr error
		code   int
	}{
		{
			&apperrors.UnauthorizedError{},
			http.StatusUnauthorized,
		},
		{
			&apperrors.QueryError{},
			http.StatusInternalServerError,
		},
		{
			&apperrors.ContextError{},
			http.StatusInternalServerError,
		},
		{
			errors.New("unknown error"),
			http.StatusInternalServerError,
		},
	}
	for _, t := range responseTests {
		s.Run(fmt.Sprintf("%T returns %d code", t.appErr, t.code), func() {
			writer := httptest.NewRecorder()

			s.base.WriteErrorResponse(ctx, writer, t.appErr)

			s.Equal(t.code, writer.Code)
			s.Equal("application/json", writer.Header().Get("Content-Type"))
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
