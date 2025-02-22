package usecases

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
)

type EnvironmentUseCase interface {
	Create(ctx context.Context, newEnvironment *entities.Environment) error
}
