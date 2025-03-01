package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
)

type EnvironmentServiceRepositoryPort interface {
	FindByProjectAndService(ctx context.Context, projectID, serviceID int) ([]*entities.EnvironmentService, error)
	DecrementAvailableRequest(ctx context.Context, environmentID, serviceID int) (*entities.EnvironmentService, error)
	Save(ctx context.Context, environmentService *entities.EnvironmentService) (*entities.EnvironmentService, error)
}
