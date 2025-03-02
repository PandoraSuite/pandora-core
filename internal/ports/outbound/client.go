package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type ClientPort interface {
	Save(ctx context.Context, client *entities.Client) (*entities.Client, error)
	FindAll(ctx context.Context, clientType enums.ClientType) ([]*entities.Client, error)
}
