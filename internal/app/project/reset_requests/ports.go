package resetrequests

import (
	"context"
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ProjectRepository interface {
	Exists(ctx context.Context, id int) (bool, errors.Error)
	GetServiceByID(ctx context.Context, id, serviceID int) (*entities.ProjectService, errors.Error)
	ResetProjectServiceUsage(ctx context.Context, id, serviceID int, nextReset time.Time) ([]*dto.EnvironmentServiceReset, errors.Error)
	ResetAvailableRequestsForEnvsService(ctx context.Context, id, serviceID int) ([]*dto.EnvironmentServiceReset, errors.Error)
}
