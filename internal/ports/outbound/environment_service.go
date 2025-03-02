package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
)

type EnvironmentServicePort interface {
	Save(ctx context.Context, environmentService *entities.EnvironmentService) (*entities.EnvironmentService, error)
	FindByProjectAndService(ctx context.Context, projectID, serviceID int) ([]*entities.EnvironmentService, error)
}

type EnvironmentServiceQuotaPort interface {
	DecrementAvailableRequest(ctx context.Context, environmentID, serviceID int) (*entities.EnvironmentService, error)
}
