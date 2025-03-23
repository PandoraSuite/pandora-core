package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ProjectPort interface {
	Save(ctx context.Context, project *entities.Project) *errors.Error
	FindByClient(ctx context.Context, clientID int) ([]*dto.ProjectResponse, *errors.Error)
}
