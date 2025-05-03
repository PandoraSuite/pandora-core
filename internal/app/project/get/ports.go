package get

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ProjectRepository interface {
	GetByID(ctx context.Context, id int) (*entities.Project, errors.Error)
}
