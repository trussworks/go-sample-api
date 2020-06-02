package services

import (
	"github.com/facebookgo/clock"
	"go.uber.org/zap"
)

// NewServiceFactory is a constructor a new ServiceFactory
func NewServiceFactory(logger *zap.Logger) ServiceFactory {
	return ServiceFactory{
		clock:  clock.New(),
		logger: logger,
	}
}

// ServiceFactory store common service utilities
// and creates service functions
type ServiceFactory struct {
	clock  clock.Clock
	logger *zap.Logger
}
