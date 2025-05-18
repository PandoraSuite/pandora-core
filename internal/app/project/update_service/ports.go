package updateservice

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ProjectRepository interface {
	Exists(ctx context.Context, id int) (bool, errors.Error)
	UpdateService(ctx context.Context, id, serviceID int, update *dto.ProjectServiceUpdate) (*entities.ProjectService, errors.Error)
	GetProjectServiceQuotaUsage(ctx context.Context, id, serviceID int) (*dto.QuotaUsage, errors.Error)
}

type EnvironmentRepository interface {
	ExistsServiceWithInfiniteMaxRequest(ctx context.Context, projectID, serviceID int) (bool, errors.Error)
}
