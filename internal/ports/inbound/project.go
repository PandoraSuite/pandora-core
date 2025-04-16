package inbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ProjectHTTPPort interface {
	Create(ctx context.Context, req *dto.ProjectCreate) (*dto.ProjectResponse, *errors.Error)
	Update(ctx context.Context, id int, req *dto.ProjectUpdate) (*dto.ProjectResponse, *errors.Error)
	GetByID(ctx context.Context, id int) (*dto.ProjectResponse, *errors.Error)
	AssignService(ctx context.Context, id int, req *dto.ProjectService) (*dto.ProjectServiceResponse, *errors.Error)
	RemoveService(ctx context.Context, id, serviceID int) *errors.Error
	GetEnvironments(ctx context.Context, id int) ([]*dto.EnvironmentResponse, *errors.Error)
}
