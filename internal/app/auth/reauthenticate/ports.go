package reauthenticate

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
	ValidateScopedToken(ctx context.Context, token, expectedScope string) (*dto.TokenResponse, errors.Error)
}
