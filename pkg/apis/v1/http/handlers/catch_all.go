package handlers

import (
	"net/http"

	"bin/bork/pkg/apperrors"
)

// NewCatchAllHandler is a constructor for CatchAllHandler
func NewCatchAllHandler(base HandlerBase) CatchAllHandler {
	return CatchAllHandler{
		base,
	}
}

// CatchAllHandler returns 404
type CatchAllHandler struct {
	HandlerBase
}

// Handle returns 404 on unexpected routes
func (h CatchAllHandler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.WriteErrorResponse(
			r.Context(),
			w,
			&apperrors.UnknownRouteError{
				Path: r.URL.Path,
			},
		)
	}
}
