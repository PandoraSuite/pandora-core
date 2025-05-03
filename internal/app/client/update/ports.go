package update

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ClientRepository interface {
	Update(ctx context.Context, id int, update *dto.ClientUpdate) (*entities.Client, errors.Error)
}
