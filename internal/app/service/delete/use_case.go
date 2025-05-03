package delete

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, id int) errors.Error
}

type useCase struct {
	validator validator.Validator

	serviceRepo    ServiceRepository
	projectRepo    ProjectRepository
	requestLogRepo RequestLogRepository
}

func (uc *useCase) Execute(ctx context.Context, id int) errors.Error {
	if err := uc.validateID(id); err != nil {
		return err
	}

	isAssigned, err := uc.projectRepo.ExistsServiceIn(ctx, id)
	if err != nil {
		return err
	}

	if isAssigned {
		return errors.NewEntityValidationFailed(
			"Service",
			"service cannot be deleted because it is assigned to projects",
			map[string]any{"id": id},
		)
	}

	if err := uc.serviceRepo.Delete(ctx, id); err != nil {
		return err
	}

	if err := uc.requestLogRepo.DeleteByService(ctx, id); err != nil {
		return err
	}

	return nil
}

func (uc *useCase) validateID(id int) errors.Error {
	return uc.validator.ValidateVariable(
		id,
		"id",
		"required,gt=0",
		map[string]string{
			"gt=0":     "id must be greater than 0",
			"required": "id is required",
		},
	)
}

func NewUseCase(
	validator validator.Validator,
	serviceRepo ServiceRepository,
	projectRepo ProjectRepository,
	requestLogRepo RequestLogRepository,
) UseCase {
	return &useCase{
		validator:      validator,
		serviceRepo:    serviceRepo,
		projectRepo:    projectRepo,
		requestLogRepo: requestLogRepo,
	}
}
