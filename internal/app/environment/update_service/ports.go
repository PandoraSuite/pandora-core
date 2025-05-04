package updateservice

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type EnvironmentRepository interface {
	Exists(ctx context.Context, id int) (bool, errors.Error)
	UpdateService(ctx context.Context, id, serviceID int, update *dto.EnvironmentServiceUpdate) (*entities.EnvironmentService, errors.Error)
	GetServiceByID(ctx context.Context, id, serviceID int) (*entities.EnvironmentService, errors.Error)
	GetProjectServiceQuotaUsage(ctx context.Context, id, serviceID int) (*dto.QuotaUsage, errors.Error)
}
