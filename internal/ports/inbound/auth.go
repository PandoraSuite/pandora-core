package inbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
)

type AuthHTTPPort interface {
	Authenticate(ctx context.Context, req *dto.AuthenticateRequest) (*dto.AuthenticateResponse, error)
	ValidateToken(ctx context.Context, req *dto.TokenRequest) (string, error)
}
