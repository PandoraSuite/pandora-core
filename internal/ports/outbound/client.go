package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
)

type ClientPort interface {
	Save(ctx context.Context, client *entities.Client) error
	FindAll(ctx context.Context, filter *dto.ClientFilter) ([]*entities.Client, error)
}
