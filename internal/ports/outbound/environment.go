package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type EnvironmentPort interface {
	Save(ctx context.Context, environment *entities.Environment) *errors.Error
	FindByID(ctx context.Context, id int) (*entities.Environment, *errors.Error)
	AddService(ctx context.Context, id int, service *entities.EnvironmentService) *errors.Error
	FindByProject(ctx context.Context, projectID int) ([]*entities.Environment, *errors.Error)
	DecrementAvailableRequest(ctx context.Context, id, serviceID int) (*dto.DecrementAvailableRequest, *errors.Error)
	ExistsEnvironmentService(ctx context.Context, id, serviceID int) (bool, *errors.Error)
	GetMaxRequestForServiceInProject(ctx context.Context, id, serviceID int) (int, *errors.Error)
	GetAllMaxRequestForServiceInEnvironments(ctx context.Context, id, serviceID int) ([]int, *errors.Error)
}
