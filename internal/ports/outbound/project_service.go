package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
)

type ProjectServiceRepositoryPort interface {
	FindByProjectAndService(ctx context.Context, projectID, serviceID int) (*entities.ProjectService, error)
	Save(ctx context.Context, projectService *entities.ProjectService) (*entities.ProjectService, error)
}
