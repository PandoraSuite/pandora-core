package usecases

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
)

type ServiceUseCase interface {
	Create(ctx context.Context, newService *entities.Service) error
}
