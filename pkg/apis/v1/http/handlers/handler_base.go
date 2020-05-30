package handlers

import (
	"github.com/facebookgo/clock"
	"go.uber.org/zap"
)

// NewHandlerBase is a constructor for HandlerBase
func NewHandlerBase(logger *zap.Logger) HandlerBase {
	return HandlerBase{
		logger: logger,
		clock:  clock.New(),
	}
}

// HandlerBase is for shared handler utilities
type HandlerBase struct {
	logger *zap.Logger
	clock  clock.Clock
}
