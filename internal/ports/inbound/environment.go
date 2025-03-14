package inbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type EnvironmentHTTPPort interface {
	Create(ctx context.Context, req *dto.EnvironmentCreate) (*dto.EnvironmentResponse, *errors.Error)
	AssignService(ctx context.Context, req *dto.AssignServiceToEnvironment) *errors.Error
	GetEnvironmentsByProject(ctx context.Context, projectID int) ([]*dto.EnvironmentResponse, *errors.Error)
}
