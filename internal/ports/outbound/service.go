package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
)

type ServicePort interface {
	Save(ctx context.Context, service *entities.Service) (*entities.Service, error)
	FindAll(ctx context.Context) ([]*entities.Service, error)
	FindActiveServices(ctx context.Context) ([]*entities.Service, error)
}

type ServiceFindPort interface {
	FindByNameAndVersion(ctx context.Context, name, version string) (*entities.Service, error)
}
