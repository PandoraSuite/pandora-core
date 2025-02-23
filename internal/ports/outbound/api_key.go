package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
)

type APIKeyRepositoryPort interface {
	Exists(ctx context.Context, key string) (bool, error)
	Save(ctx context.Context, service *entities.APIKey) (*entities.APIKey, error)
}
