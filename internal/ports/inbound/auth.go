package inbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
)

type AuthHTTPPort interface {
	Authenticate(ctx context.Context, req *dto.Authenticate) (*dto.AuthenticateResponse, error)
	ValidateToken(ctx context.Context, req *dto.TokenRequest) (string, error)
	ChangePassword(ctx context.Context, req *dto.ChangePassword) error
}
