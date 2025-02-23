package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
)

type ProjectServiceRepositoryPort interface {
	Save(ctx context.Context, service *entities.ProjectService) (*entities.ProjectService, error)
}
