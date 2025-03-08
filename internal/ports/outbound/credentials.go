package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
)

type CredentialsPort interface {
	FindCredentials(ctx context.Context, username string) (*entities.Credential, error)
}
