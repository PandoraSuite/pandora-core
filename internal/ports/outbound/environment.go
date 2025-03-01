package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
)

type EnvironmentRepositoryPort interface {
	FindByID(ctx context.Context, id int) (*entities.Environment, error)
	FindByProject(ctx context.Context, projectID int) ([]*entities.Environment, error)
	Save(ctx context.Context, environment *entities.Environment) (*entities.Environment, error)
}
