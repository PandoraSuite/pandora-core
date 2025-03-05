package inbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
)

type APIKeyHTTPPort interface {
	Create(ctx context.Context, req *dto.APIKeyCreate) (*dto.APIKeyResponse, error)
	GetAPIKeysByEnvironment(ctx context.Context, environmentID int) ([]*dto.APIKeyResponse, error)
}

type APIKeyGRPCPort interface {
	ValidateAndConsume(ctx context.Context, req *dto.APIKeyValidateAndConsume) (*dto.APIKeyValidateResponse, error)
}
