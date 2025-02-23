package inbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
)

type ProjectServicePort interface {
	Create(ctx context.Context, req *dto.ProjectServiceCreate) (*dto.ProjectServiceResponse, error)
}
