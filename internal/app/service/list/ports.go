package list

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ServiceRepository interface {
	List(ctx context.Context, filter *dto.ServiceFilter) ([]*entities.Service, errors.Error)
}
