package usecases

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
)

type RequestLogUseCase interface {
	Create(ctx context.Context, newRequestLog *entities.RequestLog) error
}
