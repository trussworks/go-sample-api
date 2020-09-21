package httpserver

import (
	"net/url"

	"bin/bork/pkg/appconfig"
)

// NewURL returns a new url for the server
func NewURL(appConfig *appconfig.AppConfig) url.URL {
	scheme := appConfig.APIProtocol
	host := appConfig.APIHost
	if appConfig.APIPort != "" {
		host = host + ":" + appConfig.APIPort
	}
	return url.URL{
		Scheme: scheme,
		Host:   host,
	}
}
