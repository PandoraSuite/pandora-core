package create

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type APIKeyRepository interface {
	Exists(ctx context.Context, key string) (bool, errors.Error)
	Create(ctx context.Context, apiKey *entities.APIKey) errors.Error
}
