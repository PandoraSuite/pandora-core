package service

import (
	"github.com/MAD-py/pandora-core/internal/app/service/create"
	"github.com/MAD-py/pandora-core/internal/app/service/delete"
	"github.com/MAD-py/pandora-core/internal/app/service/list"
	listrequest "github.com/MAD-py/pandora-core/internal/app/service/list_request"
	updatestatus "github.com/MAD-py/pandora-core/internal/app/service/update_status"
	"github.com/MAD-py/pandora-core/internal/validator"
)

// ... Create Use Case ...

type CreateUseCase = create.UseCase

func NewCreateUseCase(
	validator validator.Validator,
	serviceRepo ServiceCreateRepository,
) CreateUseCase {
	return create.NewUseCase(validator, serviceRepo)
}

// ... Delete Use Case ...

type DeleteUseCase = delete.UseCase

func NewDeleteUseCase(
	validator validator.Validator,
	serviceRepo ServiceDeleteRepository,
	projectRepo ProjectServiceVerifier,
) DeleteUseCase {
	return delete.NewUseCase(validator, serviceRepo, projectRepo)
}

// ... List Use Case ...

type ListUseCase = list.UseCase

func NewListUseCase(
	validator validator.Validator,
	serviceRepo ServiceListRepository,
) ListUseCase {
	return list.NewUseCase(validator, serviceRepo)
}

// ... List Requests Use Case ...
type ListRequestsUseCase = listrequest.UseCase

func NewListRequestsUseCase(
	validator validator.Validator,
	serviceRepo ServiceListRequestsRepository,
	requestRepo RequestListByServiceRepository,
) ListRequestsUseCase {
	return listrequest.NewUseCase(validator, serviceRepo, requestRepo)
}

// ... Update Status Use Case ...

type UpdateStatusUseCase = updatestatus.UseCase

func NewUpdateStatusUseCase(
	validator validator.Validator,
	serviceRepo ServiceUpdateStatusRepository,
) UpdateStatusUseCase {
	return updatestatus.NewUseCase(validator, serviceRepo)
}
