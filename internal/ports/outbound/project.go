package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
)

type ProjectRepositoryPort interface {
	FindByClient(ctx context.Context, clientID int) ([]*entities.Project, error)
	Save(ctx context.Context, project *entities.Project) (*entities.Project, error)
}
