package inbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
)

type EnvironmentPort interface {
	Create(ctx context.Context, req *dto.EnvironmentCreate) (*dto.EnvironmentResponse, error)
	AssignService(ctx context.Context, req *dto.AssignServiceToEnvironment) error
	GetEnvironmentsByProject(ctx context.Context, projectID int) ([]*dto.EnvironmentResponse, error)
}
