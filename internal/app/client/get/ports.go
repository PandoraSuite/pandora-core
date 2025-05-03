package get

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ClientRepository interface {
	GetByID(ctx context.Context, id int) (*entities.Client, errors.Error)
}
