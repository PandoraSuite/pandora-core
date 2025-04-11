package app

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/ports/outbound"
)

type ReservationUseCase struct {
	reservationRepo outbound.ReservationPort
	environmentRepo outbound.EnvironmentPort
}

func (u *ReservationUseCase) Commit(
	ctx context.Context, id string,
) *errors.Error {
	_, err := u.reservationRepo.RemoveReservation(ctx, id)
	return err
}

func (u *ReservationUseCase) Rollback(
	ctx context.Context, id string,
) *errors.Error {
	reservation, err := u.reservationRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if _, err := u.reservationRepo.RemoveReservation(ctx, id); err != nil {
		return err
	}
	// Must not exist reservation for increasing available request
	if err := u.environmentRepo.IncreaseAvailableRequest(
		ctx, reservation.EnvironmentID, reservation.ServiceID); err != nil {
		return err
	}
	return nil
}

func NewReservationUseCase(ReservationRepo outbound.ReservationPort,
	EnvironmentRepo outbound.EnvironmentPort) *ReservationUseCase {
	return &ReservationUseCase{
		reservationRepo: ReservationRepo,
		environmentRepo: EnvironmentRepo,
	}
}
