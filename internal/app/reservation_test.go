package app

import (
	"context"
	"testing"

	"github.com/MAD-py/pandora-core/internal/ports/outbound/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ReservationSuite struct {
	suite.Suite

	ctrl *gomock.Controller

	reservationRepo *mock.MockReservationPort
	environmentRepo *mock.MockEnvironmentPort

	useCase *ReservationUseCase

	ctx context.Context
}

func (s *ReservationSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.reservationRepo = mock.NewMockReservationPort(s.ctrl)
	s.environmentRepo = mock.NewMockEnvironmentPort(s.ctrl)

	s.useCase = NewReservationUseCase(s.reservationRepo, s.environmentRepo)

	s.ctx = context.Background()
}

func (s *ReservationSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestReservationSuite(t *testing.T) {
	suite.Run(t, new(ReservationSuite))
}
