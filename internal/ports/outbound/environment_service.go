package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
)

type EnvironmentServiceRepositoryPort interface {
	Save(ctx context.Context, service *entities.EnvironmentService) (*entities.EnvironmentService, error)
}
