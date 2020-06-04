package httpserver

import (
	"fmt"
	"net/url"

	"bin/bork/pkg/appconfig"
	"bin/bork/pkg/sources/postgres"
)

const configMissingMessage = "Must set config: %v"

func (s Server) checkRequiredConfig(config string) {
	if s.Config.GetString(config) == "" {
		s.logger.Fatal(fmt.Sprintf(configMissingMessage, config))
	}
}

// NewURL returns a new url for the server
func (s Server) NewURL() url.URL {
	s.checkRequiredConfig(appconfig.APIHostKey)
	scheme := s.Config.GetString(appconfig.APIProtocolKey)
	host := s.Config.GetString(appconfig.APIHostKey)
	if port := s.Config.GetString(appconfig.APIPortKey); port != "" {
		host = host + ":" + port
	}
	return url.URL{
		Scheme: scheme,
		Host:   host,
	}
}

// NewDBConfig returns a new DBConfig and check required fields
func (s Server) NewDBConfig() postgres.DBConfig {
	s.checkRequiredConfig(appconfig.DBHostConfigKey)
	s.checkRequiredConfig(appconfig.DBPortConfigKey)
	s.checkRequiredConfig(appconfig.DBNameConfigKey)
	s.checkRequiredConfig(appconfig.DBUsernameConfigKey)
	s.checkRequiredConfig(appconfig.DBPasswordConfigKey)
	s.checkRequiredConfig(appconfig.DBSSLModeConfigKey)
	return postgres.DBConfig{
		Host:     s.Config.GetString(appconfig.DBHostConfigKey),
		Port:     s.Config.GetString(appconfig.DBPortConfigKey),
		Database: s.Config.GetString(appconfig.DBNameConfigKey),
		Username: s.Config.GetString(appconfig.DBUsernameConfigKey),
		Password: s.Config.GetString(appconfig.DBPasswordConfigKey),
		SSLMode:  s.Config.GetString(appconfig.DBSSLModeConfigKey),
	}
}
