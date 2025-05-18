package removeservice

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type EnvironmentRepository interface {
	Exists(ctx context.Context, id int) (bool, errors.Error)
	RemoveService(ctx context.Context, id, serviceID int) (int64, errors.Error)
}
