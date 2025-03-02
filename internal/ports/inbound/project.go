package inbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
)

type ProjectPort interface {
	Create(ctx context.Context, req *dto.ProjectCreate) (*dto.ProjectResponse, error)
	AssignService(ctx context.Context, req *dto.AssignServiceToProject) error
	GetProjectsByClient(ctx context.Context, clientID int) ([]*dto.ProjectResponse, error)
}
