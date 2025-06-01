package scopedtokenvalidation

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type TokenProvider interface {
	ValidateScopedToken(ctx context.Context, token, expectedScope string) (string, errors.Error)
}
