package httpserver

import (
	"net/http"
	"time"

	"bin/bork/pkg/apis/v1/http/handlers"
	"bin/bork/pkg/appcontext"
	"bin/bork/pkg/apperrors"
	"bin/bork/pkg/models"
)

type FakeAuthorizeMiddlewareFactory struct {
	base handlers.HandlerBase
}

func (m FakeAuthorizeMiddlewareFactory) authorizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.base.WriteErrorResponse(r.Context(), w, &apperrors.UnauthorizedError{})
		} else {
			ctx := appcontext.WithUser(r.Context(), models.User{ID: authHeader})
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

// NewFakeAuthorizeMiddleware does some fake authorization
func NewFakeAuthorizeMiddleware(base handlers.HandlerBase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return FakeAuthorizeMiddlewareFactory{base: base}.authorizeMiddleware(next)
	}
}

const CookieName = "go-sample-api-session"

type SimpleSessionMiddlewareFactory struct {
	base handlers.HandlerBase
	path string
}

func (m SimpleSessionMiddlewareFactory) authorizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := appcontext.WithSessionCreator(r.Context(),
			func(session string) error {
				// The cookie settings should be configurable,
				// but as a proof of concept ...
				expires := m.base.Clock.Now().Add(20 * time.Minute)
				http.SetCookie(w, &http.Cookie{
					Name: CookieName,
					Value: session,
					Path: m.path,
					Expires: expires,
					MaxAge: 86400,
				})
				return nil
			})
		sessionCookie, err := r.Cookie(CookieName)
		if err == nil {
			// in the real world, check if the session is
			// valid and which user it is associated with
			if sessionCookie.Value != "" {
				user := models.User{ID: sessionCookie.Value}
				ctx = appcontext.WithUser(ctx, user)
			}
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NewSimpleSessionMiddleware(base handlers.HandlerBase, path string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return SimpleSessionMiddlewareFactory{base: base, path: path}.
			authorizeMiddleware(next)
	}
}
