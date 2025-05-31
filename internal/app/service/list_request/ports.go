package listrequest

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ServiceRepository interface {
	Exists(ctx context.Context, id int) (bool, errors.Error)
}

type RequestRepository interface {
	ListByService(ctx context.Context, serviceID int, filter *dto.RequestFilter) ([]*entities.Request, errors.Error)
}
