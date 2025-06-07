package delete

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type APIKeyRepository interface {
	Delete(ctx context.Context, id int) errors.Error
}
