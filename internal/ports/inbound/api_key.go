package inbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
)

type APIKeyPort interface {
	Create(ctx context.Context, req *dto.APIKeyCreate) (*dto.APIKeyResponse, error)
	ValidateAndConsume(ctx context.Context, req *dto.APIKeyValidateAndConsume) (*dto.APIKeyValidateResponse, error)
	GetAPIKeysByEnvironment(ctx context.Context, environmentID int) ([]*dto.APIKeyResponse, error)
}
