package inbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
)

type RequestLogPort interface {
	Create(ctx context.Context, req *dto.RequestLogCreate) (*dto.RequestLogResponse, error)
}
