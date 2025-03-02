package inbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
)

type ServicePort interface {
	Create(ctx context.Context, req *dto.ServiceCreate) (*dto.ServiceResponse, error)
	GetServices(ctx context.Context) ([]*dto.ServiceResponse, error)
	GetActiveServices(ctx context.Context) ([]*dto.ServiceResponse, error)
}
