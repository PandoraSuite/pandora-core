package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ReservationPort interface {
	Save(ctx context.Context, reservation *entities.Reservation) *errors.Error
	CountReservationsByFields(
		ctx context.Context, environment_id, service_id int, api_key string,
	) (int, *errors.Error)
}
