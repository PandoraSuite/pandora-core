package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
)

type APIKeyRepositoryPort interface {
	Exists(ctx context.Context, key string) (bool, error)
	Save(ctx context.Context, service *entities.APIKey) (*entities.APIKey, error)
}

type ClientRepositoryPort interface {
	Save(ctx context.Context, service *entities.Client) (*entities.Client, error)
}

type EnvironmentServiceRepositoryPort interface {
	Save(ctx context.Context, service *entities.EnvironmentService) (*entities.EnvironmentService, error)
}

type EnvironmentRepositoryPort interface {
	Save(ctx context.Context, service *entities.Environment) (*entities.Environment, error)
}

type ProjectServiceRepositoryPort interface {
	Save(ctx context.Context, service *entities.ProjectService) (*entities.ProjectService, error)
}

type ProjectRepositoryPort interface {
	Save(ctx context.Context, service *entities.Project) (*entities.Project, error)
}

type RequestLogRepositoryPort interface {
	Save(ctx context.Context, service *entities.RequestLog) (*entities.RequestLog, error)
}

type ServiceRepositoryPort interface {
	Save(ctx context.Context, service *entities.Service) (*entities.Service, error)
}
