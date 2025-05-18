package create

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type EnvironmentRepository interface {
	Create(ctx context.Context, environment *entities.Environment) errors.Error
}

type ProjectRepository interface {
	Exists(ctx context.Context, id int) (bool, errors.Error)
	GetProjectServiceQuotaUsage(ctx context.Context, id, serviceID int) (*dto.QuotaUsage, errors.Error)
}
