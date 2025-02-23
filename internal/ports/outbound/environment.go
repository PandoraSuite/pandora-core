package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
)

type EnvironmentRepositoryPort interface {
	Save(ctx context.Context, service *entities.Environment) (*entities.Environment, error)
}
