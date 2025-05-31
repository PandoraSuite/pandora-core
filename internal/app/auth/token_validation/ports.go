package tokenvalidation

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type TokenProvider interface {
	Validate(ctx context.Context, token string) (string, errors.Error)
}
