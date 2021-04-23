package appcontext

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type ContextTestSuite struct {
	suite.Suite
	logger *zap.Logger
}

func TestContextTestSuite(t *testing.T) {
	contextTestSuite := &ContextTestSuite{
		Suite:  suite.Suite{},
		logger: zap.NewNop(),
	}
	suite.Run(t, contextTestSuite)
}

func (s ContextTestSuite) TestWithTrace() {
	ctx, tID := WithTrace(context.Background())
	traceID := ctx.Value(traceKey).(uuid.UUID)

	s.NotEqual(uuid.UUID{}, traceID)
	s.Equal(tID.String(), traceID.String())
}

func (s ContextTestSuite) TestTrace() {
	ctx := context.Background()
	expectedID := uuid.New()
	ctx = context.WithValue(ctx, traceKey, expectedID)

	traceID, ok := Trace(ctx)

	s.True(ok)
	s.Equal(expectedID, traceID)
}
