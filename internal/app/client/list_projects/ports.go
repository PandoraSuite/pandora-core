package listprojects

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ClientRepository interface {
	Exists(ctx context.Context, id int) (bool, errors.Error)
}

type ProjectRepository interface {
	ListByClient(ctx context.Context, id int) ([]*entities.Project, errors.Error)
}
