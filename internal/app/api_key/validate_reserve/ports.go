package validatereserve

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type RequestRepository interface {
	Create(ctx context.Context, request *entities.RequestLog) errors.Error
}

type ReservationRepository interface {
	GetByIDWithDetails(ctx context.Context, id string) (*dto.ReservationWithDetails, errors.Error)
}
