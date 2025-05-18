package assignservice

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type EnvironmentRepository interface {
	Exists(ctx context.Context, id int) (bool, errors.Error)
	AddService(ctx context.Context, id int, service *entities.EnvironmentService) errors.Error
	ExistsServiceIn(ctx context.Context, id, serviceID int) (bool, errors.Error)
	GetProjectServiceQuotaUsage(ctx context.Context, id, serviceID int) (*dto.QuotaUsage, errors.Error)
}
