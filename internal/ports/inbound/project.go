package inbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ProjectHTTPPort interface {
	Create(ctx context.Context, req *dto.ProjectCreate) (*dto.ProjectResponse, *errors.Error)
	GetByClient(ctx context.Context, clientID int) ([]*dto.ProjectResponse, *errors.Error)
	AssignService(ctx context.Context, req *dto.AssignServiceToProject) *errors.Error
}
