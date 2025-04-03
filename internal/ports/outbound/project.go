package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ProjectPort interface {
	Save(ctx context.Context, project *entities.Project) *errors.Error
	Exists(ctx context.Context, id int) (bool, *errors.Error)
	FindByID(ctx context.Context, id int) (*entities.Project, *errors.Error)
	AddService(ctx context.Context, id int, service *entities.ProjectService) *errors.Error
	FindByClient(ctx context.Context, clientID int) ([]*entities.Project, *errors.Error)
	GetMaxRequest(ctx context.Context, id, serviceID int) (int, *errors.Error)
}
