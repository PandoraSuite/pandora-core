package validateonly

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type APIKeyRepository interface {
	GetByKey(ctx context.Context, key string) (*entities.APIKey, errors.Error)
	UpdateLastUsed(ctx context.Context, key string) errors.Error
}

type EnvironmentRepository interface {
	GetByID(ctx context.Context, id int) (*entities.Environment, errors.Error)
	ExistsServiceIn(ctx context.Context, id, serviceID int) (bool, errors.Error)
}

type ProjectRepository interface {
	GetProjectContextByID(ctx context.Context, id int) (*dto.ProjectContextResponse, errors.Error)
}

type ServiceRepository interface {
	GetByNameAndVersion(ctx context.Context, name, version string) (*entities.Service, errors.Error)
}

type RequestRepository interface {
	Create(ctx context.Context, request *entities.Request) errors.Error
}
