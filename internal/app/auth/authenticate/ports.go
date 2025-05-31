package authenticate

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type CredentialsRepository interface {
	GetByUsername(ctx context.Context, username string) (*entities.Credentials, errors.Error)
}

type TokenProvider interface {
	GenerateAccessToken(ctx context.Context, subject string) (*dto.TokenResponse, errors.Error)
}
