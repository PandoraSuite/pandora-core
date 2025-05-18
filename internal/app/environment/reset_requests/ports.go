package resetrequests

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type EnvironmentRepository interface {
	Exists(ctx context.Context, id int) (bool, errors.Error)
	ResetAvailableRequests(ctx context.Context, id, serviceID int) (*entities.EnvironmentService, errors.Error)
}
