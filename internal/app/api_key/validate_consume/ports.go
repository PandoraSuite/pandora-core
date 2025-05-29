package validateconsume

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/app/api_key/shared"
	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type APIKeyRepository interface {
	shared.ValidateAPIKeyRepository
	UpdateLastUsed(ctx context.Context, key string) errors.Error
}

type EnvironmentRepository interface {
	shared.ValidateEnvironmentRepository
	DecrementAvailableRequest(ctx context.Context, id, serviceID int) (*dto.DecrementAvailableRequest, errors.Error)
}

type ProjectRepository interface {
	shared.ValidateProjectRepository
}

type ServiceRepository interface {
	shared.ValidateServiceRepository
}

type RequestRepository interface {
	Create(ctx context.Context, request *entities.Request) errors.Error
}
