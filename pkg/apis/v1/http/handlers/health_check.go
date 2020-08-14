package handlers

import (
	"encoding/json"
	"net/http"

	"bin/bork/pkg/apperrors"
)
type HealthDatetime string
type HealthVersion string
type HealthTimestamp string

// NewHealthCheckHandler is a constructor for HealthCheckHandler
func NewHealthCheckHandler(base HandlerBase, datetime HealthDatetime, version HealthVersion, timestamp HealthTimestamp) HealthCheckHandler {
	return HealthCheckHandler{
		HandlerBase: HandlerBase{},
		datetime: datetime,
		version: version,
		timestamp: timestamp,
	}
}

// HealthCheckHandler returns the API status
type HealthCheckHandler struct {
	HandlerBase
	datetime HealthDatetime
	version  HealthVersion
	timestamp HealthTimestamp
}

type status string

const (
	statusPass status = "pass"
)

type healthCheck struct {
	Status    status `json:"status"`
	Datetime  string `json:"datetime"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
}

// Handle handles a web request and returns a health check JSON payload
func (h HealthCheckHandler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			statusReport := healthCheck{
				Status:    statusPass,
				Version:   string(h.version),
				Datetime:  string(h.datetime),
				Timestamp: string(h.timestamp),
			}
			js, err := json.Marshal(statusReport)
			if err != nil {
				h.WriteErrorResponse(r.Context(), w, err)
				return
			}

			w.Header().Set("Content-Type", "application/json")

			_, err = w.Write(js)
			if err != nil {
				h.WriteErrorResponse(r.Context(), w, err)
				return
			}
		default:
			h.WriteErrorResponse(r.Context(), w, &apperrors.MethodNotAllowedError{Method: r.Method})
			return
		}
	}
}
