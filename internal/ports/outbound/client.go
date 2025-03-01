package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type ClientRepositoryPort interface {
	FindAll(ctx context.Context, clientType enums.ClientType) ([]*entities.Client, error)
	Save(ctx context.Context, service *entities.Client) (*entities.Client, error)
}
