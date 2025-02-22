package usecases

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
)

type EnvironmentServiceUseCase interface {
	Create(ctx context.Context, newEnvironmentService *entities.EnvironmentService) error
}
