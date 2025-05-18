package rollback

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ReservationRepository interface {
	Delete(ctx context.Context, id string) errors.Error
	GetByID(ctx context.Context, id string) (*entities.Reservation, errors.Error)
}

type EnvironmentRepository interface {
	IncreaseAvailableRequest(ctx context.Context, id, serviceID int) errors.Error
}
