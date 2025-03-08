package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
)

type CredentialsPort interface {
	ChangePassword(ctx context.Context, credentials *entities.Credentials) error
	FindCredentials(ctx context.Context, username string) (*entities.Credentials, error)
}
