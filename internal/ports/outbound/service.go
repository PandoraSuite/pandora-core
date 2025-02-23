package inbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
)

type ServiceRepositoryPort interface {
	Save(ctx context.Context, service *entities.Service) (*entities.Service, error)
}
