package inbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type AuthHTTPPort interface {
	Authenticate(ctx context.Context, req *dto.Authenticate) (*dto.AuthenticateResponse, *errors.Error)
	ValidateToken(ctx context.Context, req *dto.TokenRequest) (string, *errors.Error)
	ChangePassword(ctx context.Context, req *dto.ChangePassword) *errors.Error
	IsPasswordResetRequired(ctx context.Context, username string) (bool, *errors.Error)
}
