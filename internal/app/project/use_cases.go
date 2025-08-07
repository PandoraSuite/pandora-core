package project

import (
	assignservice "github.com/MAD-py/pandora-core/internal/app/project/assign_service"
	"github.com/MAD-py/pandora-core/internal/app/project/create"
	"github.com/MAD-py/pandora-core/internal/app/project/delete"
	"github.com/MAD-py/pandora-core/internal/app/project/get"
	"github.com/MAD-py/pandora-core/internal/app/project/list"
	listenvironments "github.com/MAD-py/pandora-core/internal/app/project/list_environments"
	removeservice "github.com/MAD-py/pandora-core/internal/app/project/remove_service"
	resetduerequests "github.com/MAD-py/pandora-core/internal/app/project/reset_due_requests"
	resetrequests "github.com/MAD-py/pandora-core/internal/app/project/reset_requests"
	"github.com/MAD-py/pandora-core/internal/app/project/update"
	updateservice "github.com/MAD-py/pandora-core/internal/app/project/update_service"
	"github.com/MAD-py/pandora-core/internal/validator"
)

// ... Assign Service Use Case ...

type AssignServiceUseCase = assignservice.UseCase

func NewAssignServiceUseCase(
	validator validator.Validator, projectRepo ProjectAssignServiceRepository,
) AssignServiceUseCase {
	return assignservice.NewUseCase(validator, projectRepo)
}

// ... Create Use Case ...

type CreateUseCase = create.UseCase

func NewCreateUseCase(
	validator validator.Validator, projectRepo ProjectCreateRepository,
) CreateUseCase {
	return create.NewUseCase(validator, projectRepo)
}

// ... Delete Use Case ...

type DeleteUseCase = delete.UseCase

func NewDeleteUseCase(
	validator validator.Validator, projectRepo ProjectDeleteRepository,
) DeleteUseCase {
	return delete.NewUseCase(validator, projectRepo)
}

// ... Get Use Case ...

type GetUseCase = get.UseCase

func NewGetUseCase(
	validator validator.Validator, projectRepo ProjectGetRepository,
) GetUseCase {
	return get.NewUseCase(validator, projectRepo)
}

// ... List Use Case ...

type ListUseCase = list.UseCase

func NewListUseCase(projectRepo ProjectListRepository) ListUseCase {
	return list.NewUseCase(projectRepo)
}

// ... List Environments Use Case ...

type ListEnvironmentsUseCase = listenvironments.UseCase

func NewListEnvironmentsUseCase(
	validator validator.Validator,
	projectRepo ProjectListEnvironmentsRepository,
	environmentRepo EnvironmentListByClientRepository,
) ListEnvironmentsUseCase {
	return listenvironments.NewUseCase(validator, projectRepo, environmentRepo)
}

// ... Remove Service Use Case ...

type RemoveServiceUseCase = removeservice.UseCase

func NewRemoveServiceUseCase(
	validator validator.Validator,
	projectRepo ProjectRemoveServiceRepository,
	environmentRepo EnvironmentRemoveServiceRepository,
) RemoveServiceUseCase {
	return removeservice.NewUseCase(validator, projectRepo, environmentRepo)
}

// ... Reset Request Use Case ...

type ResetRequestUseCase = resetrequests.UseCase

func NewResetRequestUseCase(
	validator validator.Validator,
	projectRepo ProjectResetRequestRepository,
) ResetRequestUseCase {
	return resetrequests.NewUseCase(validator, projectRepo)
}

// ... Reset Due Requests Use Case ...

type ResetDueRequestsUseCase = resetduerequests.UseCase

func NewResetDueRequestsUseCase(
	projectRepo ProjectResetDueRequestsRepository,
) ResetDueRequestsUseCase {
	return resetduerequests.NewUseCase(projectRepo)
}

// ... Update Use Case ...

type UpdateUseCase = update.UseCase

func NewUpdateUseCase(
	validator validator.Validator, projectRepo ProjectUpdateRepository,
) UpdateUseCase {
	return update.NewUseCase(validator, projectRepo)
}

// ... Update Service Use Case ...

type UpdateServiceUseCase = updateservice.UseCase

func NewUpdateServiceUseCase(
	validator validator.Validator,
	projectRepo ProjectUpdateServiceRepository,
	environmentRepo EnvironmentServiceInfiniteQuotaRepository,
) UpdateServiceUseCase {
	return updateservice.NewUseCase(validator, projectRepo, environmentRepo)
}
