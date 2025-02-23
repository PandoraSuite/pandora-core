package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
)

type ProjectRepositoryPort interface {
	Save(ctx context.Context, service *entities.Project) (*entities.Project, error)
}
