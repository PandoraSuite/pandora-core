package outbound

import (
	"context"
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ProjectPort interface {
	Save(ctx context.Context, project *entities.Project) *errors.Error
	Exists(ctx context.Context, id int) (bool, *errors.Error)
	Update(ctx context.Context, id int, update *dto.ProjectUpdate) (*entities.Project, *errors.Error)
	FindByID(ctx context.Context, id int) (*entities.Project, *errors.Error)
	AddService(ctx context.Context, id int, service *entities.ProjectService) *errors.Error
	FindByClient(ctx context.Context, clientID int) ([]*entities.Project, *errors.Error)
	UpdateService(ctx context.Context, id, serviceID int, update *dto.ProjectServiceUpdate) (*entities.ProjectService, *errors.Error)
	RemoveService(ctx context.Context, id, serviceID int) (int64, *errors.Error)
	FindContextByID(ctx context.Context, id int) (*dto.ProjectContext, *errors.Error)
	FindServiceByID(ctx context.Context, id, serviceID int) (*entities.ProjectService, *errors.Error)
	ExistsServiceIn(ctx context.Context, serviceID int) (bool, *errors.Error)
	ResetProjectServiceUsage(ctx context.Context, id, serviceID int, nextReset time.Time) ([]*dto.EnvironmentServiceReset, *errors.Error)
	GetProjectServiceQuotaUsage(ctx context.Context, id, serviceID int) (*dto.QuotaUsage, *errors.Error)
	ResetAvailableRequestsForEnvsService(ctx context.Context, id, serviceID int) ([]*dto.EnvironmentServiceReset, *errors.Error)
}
