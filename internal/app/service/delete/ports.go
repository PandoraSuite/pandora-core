package delete

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ServiceRepository interface {
	Delete(ctx context.Context, id int) errors.Error
}

type ProjectRepository interface {
	ExistsServiceIn(ctx context.Context, serviceID int) (bool, errors.Error)
}

type RequestLogRepository interface {
	DeleteByService(ctx context.Context, serviceID int) errors.Error
}
