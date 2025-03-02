package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
)

type APIKeyPort interface {
	Save(ctx context.Context, apiKey *entities.APIKey) (*entities.APIKey, error)
	Exists(ctx context.Context, key string) (bool, error)
	FindByKey(ctx context.Context, key string) (*entities.APIKey, error)
	FindByEnvironment(ctx context.Context, environmentID int) ([]*entities.APIKey, error)
}
