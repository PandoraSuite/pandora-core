package app

import (
	"context"
	"testing"

	"github.com/MAD-py/pandora-core/internal/ports/outbound/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type RequestLogSuite struct {
	suite.Suite

	ctrl *gomock.Controller

	requestLogRepo *mock.MockRequestLogPort

	useCase *RequestLogUseCase

	ctx context.Context
}

func (s *RequestLogSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.requestLogRepo = mock.NewMockRequestLogPort(s.ctrl)

	s.useCase = NewRequestLogUseCase(s.requestLogRepo)

	s.ctx = context.Background()
}

func (s *RequestLogSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestRequestLogSuite(t *testing.T) {
	suite.Run(t, new(RequestLogSuite))
}
