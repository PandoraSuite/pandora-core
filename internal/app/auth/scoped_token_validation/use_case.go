package scopedtokenvalidation

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, token string, expectedScope enums.Scope) (string, errors.Error)
}

type useCase struct {
	validator validator.Validator

	tokenProvider TokenProvider
}

func (uc *useCase) Execute(
	ctx context.Context, token string, expectedScope enums.Scope,
) (string, errors.Error) {
	if err := uc.validateInput(token, expectedScope); err != nil {
		return "", err
	}

	subject, err := uc.tokenProvider.ValidateScopedToken(
		ctx, token, string(expectedScope),
	)
	if err != nil {
		return "", err
	}

	return subject, nil
}

func (uc *useCase) validateInput(token string, expectedScope enums.Scope) errors.Error {
	var err errors.Error

	if errID := uc.validateScopedToken(token); errID != nil {
		err = errors.Aggregate(err, errID)
	}

	if errReq := uc.validateExpectedScope(expectedScope); errReq != nil {
		err = errors.Aggregate(err, errReq)
	}

	return err
}

func (uc *useCase) validateScopedToken(token string) errors.Error {
	return uc.validator.ValidateVariable(
		token,
		"scoped_token",
		"required,jwt",
		map[string]string{
			"jwt":      "scoped_token must be a valid JWT",
			"required": "scoped_token is required",
		},
	)
}

func (uc *useCase) validateExpectedScope(expectedScope enums.Scope) errors.Error {
	return uc.validator.ValidateVariable(
		expectedScope,
		"expected_scope",
		"required,enums=api_key:reveal",
		map[string]string{
			"enums":    "status must be one of the following: api_key:reveal",
			"required": "expected_scope is required",
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
