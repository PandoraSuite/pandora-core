package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type EnvironmentServicePort interface {
	Save(ctx context.Context, environmentService *entities.EnvironmentService) *errors.Error
	FindByProjectAndService(ctx context.Context, projectID, serviceID int) ([]*entities.EnvironmentService, *errors.Error)
}

type EnvironmentServiceQuotaPort interface {
	DecrementAvailableRequest(ctx context.Context, environmentID, serviceID int) (*entities.EnvironmentService, *errors.Error)
}
