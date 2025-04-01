package inbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ServiceHTTPPort interface {
	Create(ctx context.Context, req *dto.ServiceCreate) (*dto.ServiceResponse, *errors.Error)
	GetServices(ctx context.Context, req *dto.ServiceFilter) ([]*dto.ServiceResponse, *errors.Error)
}
