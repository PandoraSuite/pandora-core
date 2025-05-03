package tokenvalidation

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type TokenProvider interface {
	Validate(ctx context.Context, token *dto.TokenValidation) (string, errors.Error)
}
