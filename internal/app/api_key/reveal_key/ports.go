package revealkey

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type APIKeyRepository interface {
	GetByID(ctx context.Context, id int) (*entities.APIKey, errors.Error)
}
