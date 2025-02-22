package usecases

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
)

type ProjectUseCase interface {
	Create(ctx context.Context, newProject *entities.Project) error
}
