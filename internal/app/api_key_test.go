package app

import (
	"context"
	"testing"

	"github.com/MAD-py/pandora-core/internal/ports/outbound/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type APIKeySuite struct {
	suite.Suite

	ctrl *gomock.Controller

	apiKeyRepo      *mock.MockAPIKeyPort
	requestLog      *mock.MockRequestLogPort
	serviceRepo     *mock.MockServiceFindPort
	environmentRepo *mock.MockEnvironmentPort
	reservationRepo *mock.MockReservationPort

	useCase *APIKeyUseCase

	ctx context.Context
}

func (s *APIKeySuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.apiKeyRepo = mock.NewMockAPIKeyPort(s.ctrl)
	s.requestLog = mock.NewMockRequestLogPort(s.ctrl)
	s.serviceRepo = mock.NewMockServiceFindPort(s.ctrl)
	s.environmentRepo = mock.NewMockEnvironmentPort(s.ctrl)
	s.reservationRepo = mock.NewMockReservationPort(s.ctrl)

	s.useCase = NewAPIKeyUseCase(
		s.apiKeyRepo,
		s.requestLog,
		s.serviceRepo,
		s.environmentRepo,
		s.reservationRepo,
	)

	s.ctx = context.Background()
}

func (s *APIKeySuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestAPIKeySuite(t *testing.T) {
	suite.Run(t, new(APIKeySuite))
}
