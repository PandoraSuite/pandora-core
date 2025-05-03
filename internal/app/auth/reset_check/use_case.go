package resetcheck

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, username string) (bool, errors.Error)
}

type useCase struct {
	validator validator.Validator

	credentialsRepo CredentialsRepository
}

func (uc *useCase) Execute(
	ctx context.Context, username string,
) (bool, errors.Error) {
	if err := uc.validateUsername(username); err != nil {
		return false, err
	}

	credentials, err := uc.credentialsRepo.GetByUsername(ctx, username)
	if err != nil {
		return false, err
	}

	return credentials.ForcePasswordReset, nil
}

func (uc *useCase) validateUsername(username string) errors.Error {
	return uc.validator.ValidateVariable(
		username,
		"username",
		"required",
		map[string]string{
			"required": "username is required",
		},
	)
}

func NewUseCase(
	validator validator.Validator,
	credentialsRepo CredentialsRepository,
) UseCase {
	return &useCase{
		validator:       validator,
		credentialsRepo: credentialsRepo,
	}
}
