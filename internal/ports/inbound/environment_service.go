package inbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
)

type EnvironmentServicePort interface {
	Create(ctx context.Context, req *dto.EnvironmentServiceCreate) (*dto.EnvironmentServiceResponse, error)
}
