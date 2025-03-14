package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ProjectServicePort interface {
	Save(ctx context.Context, projectService *entities.ProjectService) *errors.Error
	BulkSave(ctx context.Context, projectServices []*entities.ProjectService) *errors.Error
}

type ProjectServiceFindPort interface {
	FindByProjectAndService(ctx context.Context, projectID, serviceID int) (*entities.ProjectService, *errors.Error)
}
