package app

import (
	"context"
	"testing"

	"github.com/MAD-py/pandora-core/internal/ports/outbound/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type EnvironmentSuite struct {
	suite.Suite

	ctrl *gomock.Controller

	clientRepo  *mock.MockEnvironmentPort
	projectRepo *mock.MockProjectPort

	useCase *EnvironmentUseCase

	ctx context.Context
}

func (s *EnvironmentSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.clientRepo = mock.NewMockEnvironmentPort(s.ctrl)
	s.projectRepo = mock.NewMockProjectPort(s.ctrl)

	s.useCase = NewEnvironmentUseCase(s.clientRepo, s.projectRepo)

	s.ctx = context.Background()
}

func (s *EnvironmentSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestEnvironmentSuite(t *testing.T) {
	suite.Run(t, new(EnvironmentSuite))
}
