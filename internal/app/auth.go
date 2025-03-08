package app

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	domainErr "github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/ports/outbound"
)

type AuthUseCase struct {
	tokenProvider   outbound.TokenPort
	credentialsRepo outbound.CredentialsPort
}

func (u *AuthUseCase) Authenticate(
	ctx context.Context, req *dto.AuthenticateRequest,
) (*dto.AuthenticateResponse, error) {
	credentials, err := u.credentialsRepo.FindCredentials(ctx, req.Username)
	if err != nil {
		if err == domainErr.ErrNotFound {
			return nil, domainErr.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := credentials.VerifyPassword(req.Password); err != nil {
		return nil, err
	}

	token, err := u.tokenProvider.GenerateToken(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	token.Username = req.Username
	return token, nil
}

func (u *AuthUseCase) ValidateToken(ctx context.Context)
