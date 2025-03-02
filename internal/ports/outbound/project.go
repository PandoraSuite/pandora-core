package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
)

type ProjectPort interface {
	Save(ctx context.Context, project *entities.Project) (*entities.Project, error)
	FindByClient(ctx context.Context, clientID int) ([]*entities.Project, error)
}
