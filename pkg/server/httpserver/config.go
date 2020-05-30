package httpserver

import (
	"fmt"
	"net/url"

	"bin/bork/pkg/appconfig"
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
