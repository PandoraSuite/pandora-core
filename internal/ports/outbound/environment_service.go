package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
)

type EnvironmentServiceRepositoryPort interface {
	DecrementAvailableRequest(ctx context.Context, environmentID, serviceID int) error
	Save(ctx context.Context, service *entities.EnvironmentService) (*entities.EnvironmentService, error)
}
