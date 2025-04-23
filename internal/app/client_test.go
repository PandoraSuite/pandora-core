package app

import (
	"context"
	"testing"

	"github.com/MAD-py/pandora-core/internal/ports/outbound/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ClientSuite struct {
	suite.Suite

	ctrl *gomock.Controller

	clientRepo  *mock.MockClientPort
	projectRepo *mock.MockProjectPort

	useCase *ClientUseCase

	ctx context.Context
}

func (s *ClientSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.clientRepo = mock.NewMockClientPort(s.ctrl)
	s.projectRepo = mock.NewMockProjectPort(s.ctrl)

	s.useCase = NewClientUseCase(s.clientRepo, s.projectRepo)

	s.ctx = context.Background()
}

func (s *ClientSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestClientSuite(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}
