package accesstokenvalidation

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, token string) (string, errors.Error)
}

type useCase struct {
	validator validator.Validator

	tokenProvider TokenProvider
}

func (uc *useCase) Execute(
	ctx context.Context, token string,
) (string, errors.Error) {
	if err := uc.validateAccessToken(token); err != nil {
		return "", err
	}

	subject, err := uc.tokenProvider.ValidateAccessToken(ctx, token)
	if err != nil {
		return "", err
	}

	return subject, nil
}

func (uc *useCase) validateAccessToken(token string) errors.Error {
	return uc.validator.ValidateVariable(
		token,
		"access_token",
		"required,jwt",
		map[string]string{
			"jwt":      "access_token must be a valid JWT",
			"required": "access_token is required",
		},
	)
}

func NewUseCase(
	validator validator.Validator, tokenProvider TokenProvider,
) UseCase {
	return &useCase{
		validator:     validator,
		tokenProvider: tokenProvider,
	}
}
