package resetduerequests

import (
	"context"
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ProjectRepository interface {
	ResetProjectServiceUsage(ctx context.Context, id, serviceID int, nextReset time.Time) ([]*dto.EnvironmentServiceReset, errors.Error)
	ListProjectServiceDueForReset(ctx context.Context, today time.Time) ([]*entities.Project, errors.Error)
}
