package inbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ServiceHTTPPort interface {
	Create(ctx context.Context, req *dto.ServiceCreate) (*dto.ServiceResponse, *errors.Error)
	GetServices(ctx context.Context, req *dto.ServiceFilter) ([]*dto.ServiceResponse, *errors.Error)
	UpdateStatus(ctx context.Context, id int, status enums.ServiceStatus) (*dto.ServiceResponse, *errors.Error)
}
