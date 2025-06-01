package ports

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type TokenProvider interface {
	// ... Generate ...
	GenerateAccessToken(ctx context.Context, subject string) (*dto.TokenResponse, errors.Error)
	GenerateScopedToken(ctx context.Context, subject, scope string) (*dto.TokenResponse, errors.Error)

	// ... Validate ...
	ValidateAccessToken(ctx context.Context, token string) (string, errors.Error)
	ValidateScopedToken(ctx context.Context, token, expectedScope string) (string, errors.Error)
}
