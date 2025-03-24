package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ServicePort interface {
	Save(ctx context.Context, service *entities.Service) *errors.Error
	FindAll(ctx context.Context, filter *dto.ServiceFilter) ([]*entities.Service, *errors.Error)
}

type ServiceFindPort interface {
	FindByNameAndVersion(ctx context.Context, name, version string) (*entities.Service, *errors.Error)
}
