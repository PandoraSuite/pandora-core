package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type TokenPort interface {
	GenerateToken(ctx context.Context, subject string) (*dto.TokenResponse, *errors.Error)
	ValidateToken(ctx context.Context, token *dto.TokenRequest) (string, *errors.Error)
}
