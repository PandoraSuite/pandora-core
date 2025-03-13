package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
)

type ProjectServicePort interface {
	Save(ctx context.Context, projectService *entities.ProjectService) error
	BulkSave(ctx context.Context, projectServices []*entities.ProjectService) error
}

type ProjectServiceFindPort interface {
	FindByProjectAndService(ctx context.Context, projectID, serviceID int) (*entities.ProjectService, error)
}
