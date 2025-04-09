package app

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/ports/outbound"
)

type ReservationUseCase struct {
	reservationRepo outbound.ReservationPort
}

func (u *ReservationUseCase) Commit(
	ctx context.Context, id string,
) *errors.Error {
	_, err := u.reservationRepo.RemoveReservation(ctx, id)
	return err
}

func NewReservationUseCase(ReservationRepo outbound.ReservationPort) *ReservationUseCase {
	return &ReservationUseCase{reservationRepo: ReservationRepo}
}
