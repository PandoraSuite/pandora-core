package inbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
)

type APIKeyPort interface {
	ValidateAndConsume(ctx context.Context, req *dto.APIKeyValidateAndConsume) *dto.APIKeyValidateResponse
	Create(ctx context.Context, req *dto.APIKeyCreate) (*dto.APIKeyResponse, error)
}
