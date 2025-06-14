package disable

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, id int) errors.Error
}

type useCase struct {
	validator validator.Validator

	apikeyRepo APIKeyRepository
}

func (uc *useCase) Execute(ctx context.Context, id int) errors.Error {
	if err := uc.validateID(id); err != nil {
		return err
	}

	apiKey, err := uc.apikeyRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if !apiKey.IsEnabled() {
		return nil
	}

	if apiKey.IsExpired() {
		return errors.NewEntityValidationFailed(
			"APIKey",
			"API key is expired and cannot be disabled",
			map[string]any{"id": id},
			nil,
		)
	}

	return uc.apikeyRepo.UpdateStatus(ctx, id, enums.APIKeyStatusDisabled)
}

func (uc *useCase) validateID(id int) errors.Error {
	return uc.validator.ValidateVariable(
		id,
		"id",
		"required,gt=0",
		map[string]string{
			"gt":       "id must be greater than 0",
			"required": "id is required",
		},
	)
}

func NewUseCase(
	validator validator.Validator,
	apikeyRepo APIKeyRepository,
) UseCase {
	return &useCase{
		validator:  validator,
		apikeyRepo: apikeyRepo,
	}
}
