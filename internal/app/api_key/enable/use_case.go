package enable

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

	apiKeyRepo      APIKeyRepository
	environmentRepo EnvironmentRepository
}

func (uc *useCase) Execute(ctx context.Context, id int) errors.Error {
	if err := uc.validateID(id); err != nil {
		return err
	}

	apiKey, err := uc.apiKeyRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if apiKey.IsEnabled() {
		return nil
	}

	if apiKey.IsExpired() {
		return errors.NewEntityValidationFailed(
			"APIKey",
			"API key is expired and cannot be enabled",
			map[string]any{"id": id},
			nil,
		)
	}

	environment, err := uc.environmentRepo.GetByID(ctx, apiKey.EnvironmentID)
	if err != nil {
		return err
	}

	if !environment.IsEnabled() {
		return errors.NewEntityValidationFailed(
			"Environment",
			"Environment is disabled and cannot enable API key for it",
			map[string]any{"id": apiKey.EnvironmentID},
			nil,
		)
	}

	return uc.apiKeyRepo.UpdateStatus(ctx, id, enums.APIKeyStatusEnabled)
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
	apiKeyRepo APIKeyRepository,
	environmentRepo EnvironmentRepository,
) UseCase {
	return &useCase{
		validator:       validator,
		apiKeyRepo:      apiKeyRepo,
		environmentRepo: environmentRepo,
	}
}
