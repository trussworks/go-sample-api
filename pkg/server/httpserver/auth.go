package httpserver

import (
	"net/http"

	"bin/bork/pkg/apis/v1/http/handlers"
	"bin/bork/pkg/appcontext"
	"bin/bork/pkg/apperrors"
	"bin/bork/pkg/models"

	"go.uber.org/zap"
)

type FakeAuthorizeMiddlewareFactory struct {
	base handlers.HandlerBase
}

func (m FakeAuthorizeMiddlewareFactory) authorizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.base.WriteErrorResponse(r.Context(), w, &apperrors.UnauthorizedError{})
			return
		}
		ctx := appcontext.WithUser(r.Context(), models.User{ID: authHeader})
		appcontext.LogRequestField(ctx, zap.String("user", authHeader))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// NewFakeAuthorizeMiddleware does some fake authorization
func NewFakeAuthorizeMiddleware(base handlers.HandlerBase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return FakeAuthorizeMiddlewareFactory{base: base}.authorizeMiddleware(next)
	}
}
