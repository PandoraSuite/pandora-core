package list

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ClientRepository interface {
	List(ctx context.Context, filter *dto.ClientFilter) ([]*entities.Client, errors.Error)
}
