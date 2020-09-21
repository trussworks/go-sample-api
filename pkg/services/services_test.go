package services

import (
	"testing"

	"github.com/facebookgo/clock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type ServicesTestSuite struct {
	suite.Suite
	ServiceFactory ServiceFactory
}

func TestServicesTestSuite(t *testing.T) {
	logger := zap.NewNop()
	serviceFactory := NewServiceFactory(logger, clock.NewMock())

	servicesTestSuite := &ServicesTestSuite{
		Suite:          suite.Suite{},
		ServiceFactory: serviceFactory,
	}
	suite.Run(t, servicesTestSuite)
}
