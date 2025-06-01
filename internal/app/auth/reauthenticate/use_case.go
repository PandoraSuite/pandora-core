package reauthenticate

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, req *dto.Reauthenticate) (*dto.ReauthenticateResponse, errors.Error)
}

type useCase struct {
	validator validator.Validator

	tokenProvider   TokenProvider
	credentialsRepo CredentialsRepository
}

func (uc *useCase) Execute(
	ctx context.Context, req *dto.Reauthenticate,
) (*dto.ReauthenticateResponse, errors.Error) {
	if err := uc.validateReq(req); err != nil {
		return nil, err
	}

	credentials, err := uc.credentialsRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if err.Code() == errors.CodeNotFound {
			return nil, errors.NewUnauthorized("Invalid password", err)
		}
		return nil, err
	}

	if err := credentials.VerifyPassword(req.Password); err != nil {
		return nil, errors.NewUnauthorized("Invalid password", err.Unwrap())
	}

	scope, err := uc.actionToScope(req.Action)
	if err != nil {
		return nil, err
	}

	token, err := uc.tokenProvider.GenerateScopedToken(
		ctx, req.Username, string(scope),
	)
	if err != nil {
		return nil, err
	}

	return &dto.ReauthenticateResponse{
		TokenResponse: token,
	}, nil
}

func (uc *useCase) actionToScope(
	action enums.SensitiveAction,
) (enums.Scope, errors.Error) {
	switch action {
	case enums.SensitiveActionRevealAPIKey:
		return enums.ScopeRevealAPIKey, nil
	default:
		return enums.ScopeNull, errors.NewInternal("Invalid action", nil)
	}
}

func (uc *useCase) validateReq(req *dto.Reauthenticate) errors.Error {
	return uc.validator.ValidateStruct(
		req,
		map[string]string{
			"username.required": "username is required",
			"password.required": "password is required",
			"action.required":   "action is required",
		},
	)
}

func NewUseCase(
	validator validator.Validator,
	tokenProvider TokenProvider,
	credentialsRepo CredentialsRepository,
) UseCase {
	return &useCase{
		validator:       validator,
		tokenProvider:   tokenProvider,
		credentialsRepo: credentialsRepo,
	}
}
