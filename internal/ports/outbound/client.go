package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ClientPort interface {
	Save(ctx context.Context, client *entities.Client) *errors.Error
	FindAll(ctx context.Context, filter *dto.ClientFilter) ([]*entities.Client, *errors.Error)
}
