package app

import (
	"context"
	"testing"

	"github.com/MAD-py/pandora-core/internal/domain/errors"
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

func (s *ReservationSuite) TestCommit_Success() {
	id := "test-id"

	s.reservationRepo.EXPECT().
		Delete(s.ctx, id).
		Return(nil).
		Times(1)

	err := s.useCase.Commit(s.ctx, id)

	s.Require().Nil(err)
}

func (s *ReservationSuite) TestCommit_ReservationRepoErrors() {
	id := "test-id"

	s.reservationRepo.EXPECT().
		Delete(s.ctx, id).
		Return(errors.ErrPersistence).
		Times(1)

	err := s.useCase.Commit(s.ctx, id)

	s.Equal(errors.ErrPersistence, err)
}

func TestReservationSuite(t *testing.T) {
	suite.Run(t, new(ReservationSuite))
}
