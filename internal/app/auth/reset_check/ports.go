package resetcheck

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type CredentialsRepository interface {
	GetByUsername(ctx context.Context, username string) (*entities.Credentials, errors.Error)
}
