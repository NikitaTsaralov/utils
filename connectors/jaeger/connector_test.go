package jaeger

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TraceSuite struct {
	suite.Suite
}

func (s *TraceSuite) TestStartStop() {
	trace, err := Start(Config{
		URL:         "http://localhost:14268/api/traces",
		ServiceName: "test",
	})
	s.Require().Nil(err)

	defer func(trace *Trace, ctx context.Context) {
		err := trace.Stop(ctx)
		s.Require().Nil(err)
	}(trace, context.Background())

	s.Require().NotNil(trace)
}

func TestTraceSuite(t *testing.T) {
	suite.Run(t, new(TraceSuite))
}
