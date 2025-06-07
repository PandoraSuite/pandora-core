package resetrequests

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, id, serviceID int) (*dto.EnvironmentServiceResponse, errors.Error)
}

type useCase struct {
	validator validator.Validator

	environmentRepo EnvironmentRepository
}

func (uc *useCase) Execute(
	ctx context.Context, id, serviceID int,
) (*dto.EnvironmentServiceResponse, errors.Error) {
	if err := uc.validateInput(id, serviceID); err != nil {
		return nil, err
	}

	exists, err := uc.environmentRepo.Exists(ctx, id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.NewEntityNotFound(
			"Environment",
			"environment not found",
			map[string]any{"id": id},
			nil,
		)
	}

	service, err := uc.environmentRepo.ResetAvailableRequests(ctx, id, serviceID)
	if err != nil {
		if err.Code() == errors.CodeNotFound {
			return nil, errors.NewEntityNotFound(
				"Service",
				"Service not assigned to environment",
				map[string]any{"id": serviceID},
				err,
			)
		}
		return nil, err
	}

	return &dto.EnvironmentServiceResponse{
		ID:               service.ID,
		Name:             service.Name,
		Version:          service.Version,
		MaxRequests:      service.MaxRequests,
		AvailableRequest: service.AvailableRequest,
		AssignedAt:       service.AssignedAt,
	}, nil
}

func (uc *useCase) validateInput(id, serviceID int) errors.Error {
	var err errors.Error

	if errID := uc.validateID(id); errID != nil {
		err = errors.Aggregate(err, errID)
	}

	if errID := uc.validateServiceID(serviceID); errID != nil {
		err = errors.Aggregate(err, errID)
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
	validator validator.Validator, environmentRepo EnvironmentRepository,
) UseCase {
	return &useCase{
		validator:       validator,
		environmentRepo: environmentRepo,
	}
}
