package environment

import (
	assignservice "github.com/MAD-py/pandora-core/internal/app/environment/assign_service"
	"github.com/MAD-py/pandora-core/internal/app/environment/create"
	"github.com/MAD-py/pandora-core/internal/app/environment/get"
	removeservice "github.com/MAD-py/pandora-core/internal/app/environment/remove_service"
	resetrequests "github.com/MAD-py/pandora-core/internal/app/environment/reset_requests"
	"github.com/MAD-py/pandora-core/internal/app/environment/update"
	updateservice "github.com/MAD-py/pandora-core/internal/app/environment/update_service"
	"github.com/MAD-py/pandora-core/internal/validator"
)

// ... Assign Service Use Case ...

type AssignServiceUseCase = assignservice.UseCase

func NewAssignServiceUseCase(
	validator validator.Validator,
	environmentRepo assignservice.EnvironmentRepository,
) AssignServiceUseCase {
	return assignservice.NewUseCase(validator, environmentRepo)
}

// ... Create Use Case ...

type CreateUseCase = create.UseCase

func NewCreateUseCase(
	validator validator.Validator,
	projectRepo ProjectQuotaRepository,
	environmentRepo EnvironmentCreateRepository,
) CreateUseCase {
	return create.NewUseCase(validator, projectRepo, environmentRepo)
}

// ... Get Use Case ...

type GetUseCase = get.UseCase

func NewGetUseCase(
	validator validator.Validator,
	environmentRepo EnvironmentGetRepository,
) GetUseCase {
	return get.NewUseCase(validator, environmentRepo)
}

// ... Remove Service Use Case ...

type RemoveServiceUseCase = removeservice.UseCase

func NewRemoveServiceUseCase(
	validator validator.Validator,
	environmentRepo EnvironmentRemoveServiceRepository,
) RemoveServiceUseCase {
	return removeservice.NewUseCase(validator, environmentRepo)
}

// ... Reset Request Use Case ...

type ResetRequestUseCase = resetrequests.UseCase

func NewResetRequestUseCase(
	validator validator.Validator,
	environmentRepo EnvironmentResetRequestRepository,
) ResetRequestUseCase {
	return resetrequests.NewUseCase(validator, environmentRepo)
}

// ... Update Use Case ...

type UpdateUseCase = update.UseCase

func NewUpdateUseCase(
	validator validator.Validator,
	environmentRepo EnvironmentUpdateRepository,
) UpdateUseCase {
	return update.NewUseCase(validator, environmentRepo)
}

// ... Update Service Use Case ...

type UpdateServiceUseCase = updateservice.UseCase

func NewUpdateServiceUseCase(
	validator validator.Validator,
	environmentRepo EnvironmentUpdateServiceRepository,
) UpdateServiceUseCase {
	return updateservice.NewUseCase(validator, environmentRepo)
}
