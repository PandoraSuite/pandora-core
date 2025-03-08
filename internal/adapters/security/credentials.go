package security

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	domainErr "github.com/MAD-py/pandora-core/internal/domain/errors"
)

type credentials struct {
	username       string
	hashedPassword string

	forcePasswordReset bool
}

type CredentialsRepository struct {
	credentials *credentials
}

func (r *CredentialsRepository) FindCredentials(
	ctx context.Context, username string,
) (*entities.Credential, error) {
	if username != r.credentials.username {
		return nil, domainErr.ErrNotFound
	}

	return &entities.Credential{
		Username:           r.credentials.username,
		HashedPassword:     r.credentials.hashedPassword,
		ForcePasswordReset: r.credentials.forcePasswordReset,
	}, nil
}

func NewCredentialsRepository(
	username, hashedPassword string, forcePasswordReset bool,
) *CredentialsRepository {
	return &CredentialsRepository{
		credentials: &credentials{
			username:           username,
			hashedPassword:     hashedPassword,
			forcePasswordReset: forcePasswordReset,
		},
	}
}
