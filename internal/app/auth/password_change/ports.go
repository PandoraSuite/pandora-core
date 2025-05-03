package passwordchange

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type CredentialsRepository interface {
	ChangePassword(ctx context.Context, credentials *entities.Credentials) errors.Error
	GetByUsername(ctx context.Context, username string) (*entities.Credentials, errors.Error)
}
