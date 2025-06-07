package client

import (
	"github.com/MAD-py/pandora-core/internal/app/client/create"
	"github.com/MAD-py/pandora-core/internal/app/client/get"
	"github.com/MAD-py/pandora-core/internal/app/client/list"
	listprojects "github.com/MAD-py/pandora-core/internal/app/client/list_projects"
	"github.com/MAD-py/pandora-core/internal/app/client/update"
	"github.com/MAD-py/pandora-core/internal/validator"
)

// ... Create Use Case ...

type CreateUseCase = create.UseCase

func NewCreateUseCase(
	validator validator.Validator, clientRepo ClientCreateRepository,
) CreateUseCase {
	return create.NewUseCase(validator, clientRepo)
}

// ... Delete Use Case ...

type DeleteUseCase = create.UseCase

func NewDeleteUseCase(
	validator validator.Validator, clientRepo ClientDeleteRepository,
) DeleteUseCase {
	return create.NewUseCase(validator, clientRepo)
}

// ... Get Use Case ...

type GetUseCase = get.UseCase

func NewGetUseCase(
	validator validator.Validator, clientRepo ClientGetRepository,
) GetUseCase {
	return get.NewUseCase(validator, clientRepo)
}

// ... List Use Case ...

type ListUseCase = list.UseCase

func NewListUseCase(
	validator validator.Validator, clientRepo ClientListRepository,
) ListUseCase {
	return list.NewUseCase(validator, clientRepo)
}

// ... List Projects Use Case ...

type ListProjectsUseCase = listprojects.UseCase

func NewListProjectsUseCase(
	validator validator.Validator,
	clientRepo ClientListProjectsRepository,
	projectRepo ProjectListByClientRepository,
) ListProjectsUseCase {
	return listprojects.NewUseCase(validator, clientRepo, projectRepo)
}

// ... Update Use Case ...

type UpdateUseCase = update.UseCase

func NewUpdateUseCase(
	validator validator.Validator, clientRepo ClientUpdateRepository,
) UpdateUseCase {
	return update.NewUseCase(validator, clientRepo)
}
