package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
)

type TokenPort interface {
	GenerateToken(ctx context.Context, subject string) (*dto.AuthenticateResponse, error)
	ValidateToken(ctx context.Context, token *dto.TokenRequest) (string, error)
}
