package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/spf13/viper"
)

func (s HandlerTestSuite) TestHealthCheckHandler_Handle() {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/healthcheck", bytes.NewBufferString(""))
	s.NoError(err)

	mockViper := viper.New()
	mockViper.SetDefault("APPLICATION_DATETIME", "mockdatetime")
	mockViper.SetDefault("APPLICATION_TS", "mocktimestamp")
	mockViper.SetDefault("APPLICATION_VERSION", "mockversion")

	mockHealthCheckHandler := NewHealthCheckHandler(
		s.base,
		mockViper,
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
