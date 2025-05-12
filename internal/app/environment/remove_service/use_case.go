package removeservice

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, id, serviceID int) errors.Error
}

type useCase struct {
	validator validator.Validator

	environmentRepo EnvironmentRepository
}

func (uc *useCase) Execute(ctx context.Context, id, serviceID int) errors.Error {
	if err := uc.validateInput(id, serviceID); err != nil {
		return err
	}

	exists, err := uc.environmentRepo.Exists(ctx, id)
	if err != nil {
		return err
	}

	if !exists {
		return errors.NewEntityNotFound(
			"Environment",
			"environment not found",
			map[string]any{"id": id},
			nil,
		)
	}

	n, err := uc.environmentRepo.RemoveService(ctx, id, serviceID)
	if err != nil {
		return err
	}

	if n == 0 {
		return errors.NewEntityNotFound(
			"Service",
			"service not found",
			map[string]any{"id": serviceID},
			nil,
		)
	}

	return nil
}

func (uc *useCase) validateInput(id, serviceID int) errors.Error {
	var err errors.Error

	if errID := uc.validateID(id); errID != nil {
		err = errors.Aggregate(err, errID)
	}

	if errServiceID := uc.validateServiceID(serviceID); errServiceID != nil {
		err = errors.Aggregate(err, errServiceID)
	}

	return err
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

func (uc *useCase) validateServiceID(serviceID int) errors.Error {
	return uc.validator.ValidateVariable(
		serviceID,
		"service_id",
		"required,gt=0",
		map[string]string{
			"gt":       "service_id must be greater than 0",
			"required": "service_id is required",
		},
	)
}

func NewUseCase(
	validator validator.Validator,
	environmentRepo EnvironmentRepository,
) UseCase {
	return &useCase{
		validator:       validator,
		environmentRepo: environmentRepo,
	}
}
