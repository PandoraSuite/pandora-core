package inbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
)

type APIKeyPort interface {
	Create(ctx context.Context, req *dto.APIKeyCreate) (*dto.APIKeyResponse, error)
}

type ClientPort interface {
	Create(ctx context.Context, req *dto.ClientCreate) (*dto.ClientResponse, error)
}

type EnvironmentServicePort interface {
	Create(ctx context.Context, req *dto.EnvironmentServiceCreate) (*dto.EnvironmentServiceResponse, error)
}

type EnvironmentPort interface {
	Create(ctx context.Context, req *dto.EnvironmentCreate) (*dto.EnvironmentResponse, error)
}

type ProjectServicePort interface {
	Create(ctx context.Context, req *dto.ProjectServiceCreate) (*dto.ProjectServiceResponse, error)
}

type ProjectPort interface {
	Create(ctx context.Context, req *dto.ProjectCreate) (*dto.ProjectResponse, error)
}

type RequestLogPort interface {
	Create(ctx context.Context, req *dto.RequestLogCreate) (*dto.RequestLogResponse, error)
}

type ServicePort interface {
	Create(ctx context.Context, req *dto.ServiceCreate) (*dto.ServiceResponse, error)
}
