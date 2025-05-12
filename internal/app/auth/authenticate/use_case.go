package authenticate

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, req *dto.Authenticate) (*dto.AuthenticateResponse, errors.Error)
}

type useCase struct {
	validator validator.Validator

	tokenProvider   TokenProvider
	credentialsRepo CredentialsRepository
}

func (uc *useCase) Execute(
	ctx context.Context, req *dto.Authenticate,
) (*dto.AuthenticateResponse, errors.Error) {
	if err := uc.validateReq(req); err != nil {
		return nil, err
	}

	credentials, err := uc.credentialsRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if err.Code() == errors.CodeNotFound {
			return nil, errors.NewUnauthorized(
				"invalid username or password",
				err,
			)
		}
		return nil, err
	}

	if err := credentials.VerifyPassword(req.Password); err != nil {
		return nil, err
	}

	token, err := uc.tokenProvider.Generate(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	return &dto.AuthenticateResponse{
		TokenResponse:      token,
		ForcePasswordReset: credentials.ForcePasswordReset,
	}, nil
}

func (uc *useCase) validateReq(req *dto.Authenticate) errors.Error {
	return uc.validator.ValidateStruct(
		req,
		map[string]string{
			"username.required": "username is required",
			"password.required": "password is required",
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
