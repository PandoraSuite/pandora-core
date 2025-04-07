package inbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type APIKeyHTTPPort interface {
	Create(ctx context.Context, req *dto.APIKeyCreate) (*dto.APIKeyResponse, *errors.Error)
	GetAPIKeysByEnvironment(ctx context.Context, environmentID int) ([]*dto.APIKeyResponse, *errors.Error)
}

type APIKeyGRPCPort interface {
	Validate(ctx context.Context, req *dto.APIKeyValidateBooking) (*dto.APIKeyValidateResponse, *errors.Error)
	ValidateAndConsume(ctx context.Context, req *dto.APIKeyValidate) (*dto.APIKeyValidateResponse, *errors.Error)
	ValidateAndBooking(ctx context.Context, req *dto.APIKeyValidate) (*dto.APIKeyValidateBookingResponse, *errors.Error)
}
