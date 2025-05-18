package create

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ServiceRepository interface {
	Create(ctx context.Context, service *entities.Service) errors.Error
}
