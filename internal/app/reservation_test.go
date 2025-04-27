package app

import (
	"context"
	"testing"
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
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

func (s *ReservationSuite) TestRollback_Success() {
	id := "test-id"

	now := time.Now().UTC()

	mockReservation := &entities.Reservation{
		ID:             id,
		EnvironmentID:  1,
		ServiceID:      2,
		APIKey:         "test-api-key",
		StartRequestID: "test-start-request-id",
		RequestTime:    now.Add(-time.Hour),
		ExpiresAt:      now.Add(12 * time.Hour),
	}

	s.reservationRepo.EXPECT().
		FindByID(s.ctx, id).
		Return(mockReservation, nil).
		Times(1)

	s.reservationRepo.EXPECT().
		Delete(s.ctx, id).
		Return(nil).
		Times(1)

	s.environmentRepo.EXPECT().
		IncreaseAvailableRequest(
			s.ctx,
			mockReservation.EnvironmentID,
			mockReservation.ServiceID,
		).
		Return(nil).
		Times(1)

	err := s.useCase.Rollback(s.ctx, id)

	s.Require().Nil(err)
}

func (s *ReservationSuite) TestRollback_FindByIDError() {
	id := "test-id"

	s.reservationRepo.EXPECT().
		FindByID(s.ctx, id).
		Return(nil, errors.ErrPersistence).
		Times(1)

	s.reservationRepo.EXPECT().
		Delete(s.ctx, id).
		Times(0)

	s.environmentRepo.EXPECT().
		IncreaseAvailableRequest(s.ctx, gomock.Any(), gomock.Any()).
		Times(0)

	err := s.useCase.Rollback(s.ctx, id)

	s.Equal(errors.ErrPersistence, err)
}

func (s *ReservationSuite) TestRollback_DeleteError() {
	id := "test-id"

	now := time.Now().UTC()

	mockReservation := &entities.Reservation{
		ID:             id,
		EnvironmentID:  1,
		ServiceID:      2,
		APIKey:         "test-api-key",
		StartRequestID: "test-start-request-id",
		RequestTime:    now.Add(-time.Hour),
		ExpiresAt:      now.Add(12 * time.Hour),
	}

	s.reservationRepo.EXPECT().
		FindByID(s.ctx, id).
		Return(mockReservation, nil).
		Times(1)

	s.reservationRepo.EXPECT().
		Delete(s.ctx, id).
		Return(errors.ErrPersistence).
		Times(1)

	s.environmentRepo.EXPECT().
		IncreaseAvailableRequest(s.ctx, gomock.Any(), gomock.Any()).
		Times(0)

	err := s.useCase.Rollback(s.ctx, id)

	s.Equal(errors.ErrPersistence, err)
}

func (s *ReservationSuite) TestRollback_IncreaseAvailableRequestError() {
	id := "test-id"

	now := time.Now().UTC()

	mockReservation := &entities.Reservation{
		ID:             id,
		EnvironmentID:  1,
		ServiceID:      2,
		APIKey:         "test-api-key",
		StartRequestID: "test-start-request-id",
		RequestTime:    now.Add(-time.Hour),
		ExpiresAt:      now.Add(12 * time.Hour),
	}

	s.reservationRepo.EXPECT().
		FindByID(s.ctx, id).
		Return(mockReservation, nil).
		Times(1)

	s.reservationRepo.EXPECT().
		Delete(s.ctx, id).
		Return(nil).
		Times(1)

	s.environmentRepo.EXPECT().
		IncreaseAvailableRequest(
			s.ctx,
			mockReservation.EnvironmentID,
			mockReservation.ServiceID,
		).
		Return(errors.ErrPersistence).
		Times(1)

	err := s.useCase.Rollback(s.ctx, id)

	s.Equal(errors.ErrPersistence, err)
}

func TestReservationSuite(t *testing.T) {
	suite.Run(t, new(ReservationSuite))
}
