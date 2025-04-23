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
	return u.reservationRepo.Delete(ctx, id)
}

func (u *ReservationUseCase) Rollback(
	ctx context.Context, id string,
) *errors.Error {
	reservation, err := u.reservationRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := u.reservationRepo.Delete(ctx, id); err != nil {
		return err
	}
	// Must not exist reservation for increasing available request
	if err := u.environmentRepo.IncreaseAvailableRequest(
		ctx, reservation.EnvironmentID, reservation.ServiceID); err != nil {
		return err
	}
	return nil
}

func NewReservationUseCase(reservationRepo outbound.ReservationPort,
	environmentRepo outbound.EnvironmentPort) *ReservationUseCase {
	return &ReservationUseCase{
		reservationRepo: reservationRepo,
		environmentRepo: environmentRepo,
	}
}
