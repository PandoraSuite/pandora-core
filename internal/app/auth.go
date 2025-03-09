package app

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	domainErr "github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/ports/outbound"
)

type AuthUseCase struct {
	tokenProvider   outbound.TokenPort
	credentialsRepo outbound.CredentialsPort
}

func (u *AuthUseCase) Authenticate(
	ctx context.Context, req *dto.Authenticate,
) (*dto.AuthenticateResponse, error) {
	credentials, err := u.credentialsRepo.FindCredentials(ctx, req.Username)
	if err != nil {
		if err == domainErr.ErrCredentialsNotFound {
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

	return token, nil
}

func (u *AuthUseCase) ChangePassword(
	ctx context.Context, req *dto.ChangePassword,
) error {
	if len(req.NewPassword) < 12 {
		return domainErr.ErrPasswordTooShort
	}

	if req.NewPassword != req.ConfirmPassword {
		return domainErr.ErrPasswordMismatch
	}

	credentials := &entities.Credentials{Username: req.Username}
	if err := credentials.CalculatePasswordHash(req.NewPassword); err != nil {
		return err
	}

	return u.credentialsRepo.ChangePassword(ctx, credentials)
}

func (u *AuthUseCase) ValidateToken(
	ctx context.Context, req *dto.TokenRequest,
) (string, error) {
	return u.tokenProvider.ValidateToken(ctx, req)
}

func (u *AuthUseCase) IsPasswordResetRequired(
	ctx context.Context, username string,
) (bool, error) {
	credentials, err := u.credentialsRepo.FindCredentials(ctx, username)
	if err != nil {
		return false, err
	}
	return credentials.ForcePasswordReset, nil
}

func NewAuthUseCase(
	tokenProvider outbound.TokenPort,
	credentialsRepo outbound.CredentialsPort,
) *AuthUseCase {
	return &AuthUseCase{
		tokenProvider:   tokenProvider,
		credentialsRepo: credentialsRepo,
	}
}
