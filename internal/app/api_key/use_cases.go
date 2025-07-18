package apikey

import (
	"github.com/MAD-py/pandora-core/internal/app/api_key/create"
	"github.com/MAD-py/pandora-core/internal/app/api_key/delete"
	"github.com/MAD-py/pandora-core/internal/app/api_key/disable"
	"github.com/MAD-py/pandora-core/internal/app/api_key/enable"
	revealkey "github.com/MAD-py/pandora-core/internal/app/api_key/reveal_key"
	"github.com/MAD-py/pandora-core/internal/app/api_key/update"
	validateconsume "github.com/MAD-py/pandora-core/internal/app/api_key/validate_consume"
	validateonly "github.com/MAD-py/pandora-core/internal/app/api_key/validate_only"
	"github.com/MAD-py/pandora-core/internal/validator"
)

// ... Create Use Case ...

type CreateUseCase = create.UseCase

func NewCreateUseCase(
	validator validator.Validator, repo APIKeyCreateRepository,
) CreateUseCase {
	return create.NewUseCase(validator, repo)
}

// ... Delete Use Case ...

type DeleteUseCase = delete.UseCase

func NewDeleteUseCase(
	validator validator.Validator, repo APIKeyDeleteRepository,
) DeleteUseCase {
	return delete.NewUseCase(validator, repo)
}

// ... Update Use Case ...

type UpdateUseCase = update.UseCase

func NewUpdateUseCase(
	validator validator.Validator, repo APIKeyUpdateRepository,
) UpdateUseCase {
	return update.NewUseCase(validator, repo)
}

// ... Validate Use Case ...

type ValidateUseCase = validateonly.UseCase

func NewValidateUseCase(
	validator validator.Validator,
	apiKeyRepo APIKeyValidateRepository,
	projectRepo ProjectValidateRepository,
	serviceRepo ServiceValidateRepository,
	requestRepo RequestValidateRepository,
	environmentRepo EnvironmentValidateRepository,
) ValidateUseCase {
	return validateonly.NewUseCase(
		validator,
		apiKeyRepo,
		projectRepo,
		serviceRepo,
		requestRepo,
		environmentRepo,
	)
}

// ... Validate And Consume Use Case ...

type ValidateConsumeUseCase = validateconsume.UseCase

func NewValidateConsumeUseCase(
	validator validator.Validator,
	apiKeyRepo APIKeyValidateConsumeRepository,
	projectRepo ProjectValidateConsumeRepository,
	requestRepo RequestValidateConsumeRepository,
	serviceRepo ServiceValidateConsumeRepository,
	environmentRepo EnvironmentValidateConsumeRepository,
) ValidateConsumeUseCase {
	return validateconsume.NewUseCase(
		validator,
		apiKeyRepo,
		projectRepo,
		serviceRepo,
		requestRepo,
		environmentRepo,
	)
}

// ... Disable Use Case ...

type DisableUseCase = disable.UseCase

func NewDisableUseCase(
	validator validator.Validator, repo APIKeyDisableRepository,
) DisableUseCase {
	return disable.NewUseCase(validator, repo)
}

// ... Enable Use Case ...

type EnableUseCase = enable.UseCase

func NewEnableUseCase(
	validator validator.Validator,
	apiKeyRepo APIKeyEnableRepository,
	environmentRepo EnvironmentEnableRepository,
) EnableUseCase {
	return enable.NewUseCase(validator, apiKeyRepo, environmentRepo)
}

// ... Reveal Key Use Case ...

type RevealKeyUseCase = revealkey.UseCase

func NewRevealKeyUseCase(
	validator validator.Validator, repo APIKeyRevealKeyRepository,
) RevealKeyUseCase {
	return revealkey.NewUseCase(validator, repo)
}
