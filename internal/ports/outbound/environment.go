package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type EnvironmentPort interface {
	Save(ctx context.Context, environment *entities.Environment) *errors.Error
	Exists(ctx context.Context, id int) (bool, *errors.Error)
	Update(ctx context.Context, id int, update *dto.EnvironmentUpdate) (*entities.Environment, *errors.Error)
	FindByID(ctx context.Context, id int) (*entities.Environment, *errors.Error)
	IsActive(ctx context.Context, id int) (bool, *errors.Error)
	AddService(ctx context.Context, id int, service *entities.EnvironmentService) *errors.Error
	RemoveService(ctx context.Context, id, serviceID int) (int64, *errors.Error)
	FindByProject(ctx context.Context, projectID int) ([]*entities.Environment, *errors.Error)
	ExistsServiceIn(ctx context.Context, id, serviceID int) (bool, *errors.Error)
	ResetAvailableRequests(ctx context.Context, id, serviceID int) (*entities.EnvironmentService, *errors.Error)
	DecrementAvailableRequest(ctx context.Context, id, serviceID int) (*dto.DecrementAvailableRequest, *errors.Error)
	GetProjectServiceQuotaUsage(ctx context.Context, id, serviceID int) (*dto.QuotaUsage, *errors.Error)
	RemoveServiceFromProjectEnvironments(ctx context.Context, projectID, serviceID int) (int64, *errors.Error)
	IncreaseAvailableRequest(ctx context.Context, id, serviceID int) *errors.Error
}
