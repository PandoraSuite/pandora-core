package ports

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type TokenProvider interface {
	// ... Generate ...
	GenerateAccessToken(ctx context.Context, subject string) (*dto.TokenResponse, errors.Error)
	GenerateSensitiveToken(ctx context.Context, subject, scope string) (*dto.TokenResponse, errors.Error)

	// ... Validate ...
	Validate(ctx context.Context, token string) (string, errors.Error)
}
