package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type CredentialsPort interface {
	ChangePassword(ctx context.Context, credentials *entities.Credentials) *errors.Error
	FindCredentials(ctx context.Context, username string) (*entities.Credentials, *errors.Error)
}
