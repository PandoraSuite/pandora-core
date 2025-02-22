package usecases

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
)

type ProjectServiceUseCase interface {
	Create(ctx context.Context, newProjectService *entities.ProjectService) error
}
