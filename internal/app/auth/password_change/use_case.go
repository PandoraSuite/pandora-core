package passwordchange

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, req *dto.ChangePassword) errors.Error
}

type useCase struct {
	validator validator.Validator

	credentialsRepo CredentialsRepository
}

func (uc *useCase) Execute(
	ctx context.Context, req *dto.ChangePassword,
) errors.Error {
	if err := uc.validateReq(req); err != nil {
		return err
	}

	currentCredentials, err := uc.credentialsRepo.GetByUsername(
		ctx, req.Username,
	)
	if err != nil {
		if err.Code() == errors.CodeNotFound {
			return errors.NewInternal("failed to change Password", nil)
		}
		return err
	}

	if err := currentCredentials.VerifyPassword(req.NewPassword); err == nil {
		return errors.NewValidationFailed(
			"new password must be different from the old password",
		)
	}

	credentials := &entities.Credentials{Username: req.Username}
	if err := credentials.CalculatePasswordHash(req.NewPassword); err != nil {
		return err
	}

	if err := uc.credentialsRepo.ChangePassword(ctx, credentials); err != nil {
		if err.Code() == errors.CodeNotFound {
			return errors.NewInternal("failed to change Password", nil)
		}
		return err
	}
	return nil
}

func (uc *useCase) validateReq(req *dto.ChangePassword) errors.Error {
	return uc.validator.ValidateStruct(
		req,
		map[string]string{
			"new_password.min":          "new_password must be at least 12 characters long",
			"username.required":         "username is required",
			"new_password.eqfield":      "new_password must match confirm_password",
			"new_password.required":     "new_password is required",
			"confirm_password.required": "confirm_password is required",
		},
	)
}

func NewUseCase(
	validator validator.Validator, credentialsRepo CredentialsRepository,
) UseCase {
	return &useCase{
		validator:       validator,
		credentialsRepo: credentialsRepo,
	}
}
