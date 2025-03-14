package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type APIKeyPort interface {
	Save(ctx context.Context, apiKey *entities.APIKey) *errors.Error
	Exists(ctx context.Context, key string) (bool, *errors.Error)
	FindByKey(ctx context.Context, key string) (*entities.APIKey, *errors.Error)
	FindByEnvironment(ctx context.Context, environmentID int) ([]*entities.APIKey, *errors.Error)
}
