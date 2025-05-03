package tokenvalidation

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, token *dto.TokenValidation) (string, errors.Error)
}

type useCase struct {
	validator validator.Validator

	tokenProvider TokenProvider
}

func (uc *useCase) Execute(
	ctx context.Context, req *dto.TokenValidation,
) (string, errors.Error) {
	if err := uc.validateReq(req); err != nil {
		return "", err
	}

	subject, err := uc.tokenProvider.Validate(ctx, req)
	if err != nil {
		return "", err
	}

	return subject, nil
}

func (uc *useCase) validateReq(req *dto.TokenValidation) errors.Error {
	return uc.validator.ValidateStruct(
		req,
		map[string]string{
			"access_token.jwt":      "access_token must be a valid JWT",
			"token_type.required":   "token_type is required",
			"access_token.required": "access_token is required",
		},
	)
}

func NewUseCase(
	validator validator.Validator,
	tokenProvider TokenProvider,
) UseCase {
	return &useCase{
		validator:     validator,
		tokenProvider: tokenProvider,
	}
}
