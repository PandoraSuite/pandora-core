package listenvironments

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ProjectRepository interface {
	Exists(ctx context.Context, id int) (bool, errors.Error)
}

type EnvironmentRepository interface {
	ListByProject(ctx context.Context, projectID int) ([]*entities.Environment, errors.Error)
}
