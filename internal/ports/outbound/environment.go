package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type EnvironmentPort interface {
	Save(ctx context.Context, environment *entities.Environment) *errors.Error
	FindByID(ctx context.Context, id int) (*entities.Environment, *errors.Error)
	FindByProject(ctx context.Context, projectID int) ([]*entities.Environment, *errors.Error)
}
