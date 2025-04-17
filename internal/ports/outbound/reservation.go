package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ReservationPort interface {
	Save(ctx context.Context, reservation *entities.Reservation) *errors.Error
	Delete(ctx context.Context, id string) *errors.Error
	FindByID(ctx context.Context, id string) (*entities.Reservation, *errors.Error)
	FindByIDWithDetails(ctx context.Context, id string) (*dto.ReservationWithDetails, *errors.Error)
	CountByEnvironmentAndService(ctx context.Context, environment_id, service_id int) (int, *errors.Error)
}
