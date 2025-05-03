package create

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ClientRepository interface {
	Create(ctx context.Context, client *entities.Client) errors.Error
}
