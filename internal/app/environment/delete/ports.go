package delete

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type EnvironmentRepository interface {
	Delete(ctx context.Context, id int) errors.Error
}
