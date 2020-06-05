package integration

// integration is a package for testing application routes
// it should attempt to mock as few dependencies as possible
// and simulate production application use

import (
	"net/http/httptest"
	"testing"

	"github.com/facebookgo/clock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"

	"bin/bork/pkg/server/httpserver"
)

type IntegrationTestSuite struct {
	suite.Suite
	server *httptest.Server
	clock  *clock.Mock
}

func TestIntegrationTestSuite(t *testing.T) {
	config := viper.New()
	config.AutomaticEnv()

	if !testing.Short() {
		server := httpserver.NewServer(config)
		testServer := httptest.NewServer(server)
		defer testServer.Close()

		testSuite := &IntegrationTestSuite{
			Suite:  suite.Suite{},
			server: testServer,
			clock:  clock.NewMock(),
		}

		suite.Run(t, testSuite)
	}
}
