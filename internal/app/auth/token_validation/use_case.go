package tokenvalidation

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, accessToken string) (string, errors.Error)
}

type useCase struct {
	validator validator.Validator

	tokenProvider TokenProvider
}

func (uc *useCase) Execute(
	ctx context.Context, accessToken string,
) (string, errors.Error) {
	if err := uc.validateAccessToken(accessToken); err != nil {
		return "", err
	}

	subject, err := uc.tokenProvider.Validate(ctx, accessToken)
	if err != nil {
		return "", err
	}

	return subject, nil
}

func (uc *useCase) validateAccessToken(accessToken string) errors.Error {
	return uc.validator.ValidateVariable(
		accessToken,
		"access_token",
		"required,jwt",
		map[string]string{
			"access_token.jwt":      "access_token must be a valid JWT",
			"access_token.required": "access_token is required",
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
