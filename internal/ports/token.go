package ports

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type TokenProvider interface {
	// ... Generate ...
	Generate(ctx context.Context, subject string) (*dto.TokenResponse, errors.Error)

	// ... Validate ...
	Validate(ctx context.Context, token *dto.TokenValidation) (string, errors.Error)
}
