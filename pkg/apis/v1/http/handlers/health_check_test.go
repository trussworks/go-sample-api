package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func (s HandlerTestSuite) TestHealthCheckHandler_Handle() {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/healthcheck", bytes.NewBufferString(""))
	s.NoError(err)

	mockHealthCheckHandler := NewHealthCheckHandler(
		s.base,
		HealthDatetime("mockdatetime"),
		HealthVersion("mockversion"),
		HealthTimestamp("mocktimestamp"),
	)
	mockHealthCheckHandler.Handle()(rr, req)

	s.Equal(http.StatusOK, rr.Code)

	var healthCheckActual healthCheck
	err = json.Unmarshal(rr.Body.Bytes(), &healthCheckActual)

	s.NoError(err)
	s.Equal(statusPass, healthCheckActual.Status)
	s.Equal("mockdatetime", healthCheckActual.Datetime)
	s.Equal("mockversion", healthCheckActual.Version)
	s.Equal("mocktimestamp", healthCheckActual.Timestamp)
}
