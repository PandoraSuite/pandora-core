package inbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
)

type ProjectPort interface {
	Create(ctx context.Context, req *dto.ProjectCreate) (*dto.ProjectResponse, error)
}
