package inbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type EnvironmentHTTPPort interface {
	Create(ctx context.Context, req *dto.EnvironmentCreate) (*dto.EnvironmentResponse, *errors.Error)
	Update(ctx context.Context, id int, req *dto.EnvironmentUpdate) (*dto.EnvironmentResponse, *errors.Error)
	GetByID(ctx context.Context, id int) (*dto.EnvironmentResponse, *errors.Error)
	RemoveService(ctx context.Context, id, serviceID int) *errors.Error
	AssignService(ctx context.Context, id int, req *dto.EnvironmentService) (*dto.EnvironmentServiceResponse, *errors.Error)
	ResetServiceRequests(ctx context.Context, id, serviceID int) (*dto.EnvironmentServiceResponse, *errors.Error)
}
