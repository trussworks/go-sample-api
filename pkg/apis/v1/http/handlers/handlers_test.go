package handlers

import (
	"testing"

	"github.com/facebookgo/clock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type HandlerTestSuite struct {
	suite.Suite
	base HandlerBase
}

func TestHandlerTestSuite(t *testing.T) {
	handlerTestSuite := &HandlerTestSuite{
		Suite: suite.Suite{},
		base:  NewHandlerBase(zap.NewNop(), clock.NewMock()),
	}
	suite.Run(t, handlerTestSuite)
}
