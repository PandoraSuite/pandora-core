package updatestatus

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, id int, status enums.ServiceStatus) (*dto.ServiceResponse, errors.Error)
}

type useCase struct {
	validator validator.Validator

	serviceRepo ServiceRepository
}

func (uc *useCase) Execute(
	ctx context.Context, id int, status enums.ServiceStatus,
) (*dto.ServiceResponse, errors.Error) {
	if err := uc.validateInput(id, status); err != nil {
		return nil, err
	}

	service, err := uc.serviceRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		if err.Code() == errors.CodeNotFound {
			return nil, errors.NewEntityNotFound(
				"Service",
				"service not found",
				map[string]any{"id": id},
				err,
			)
		}
		return nil, err
	}

	return &dto.ServiceResponse{
		ID:        service.ID,
		Name:      service.Name,
		Status:    service.Status,
		Version:   service.Version,
		CreatedAt: service.CreatedAt,
	}, nil
}

func (uc *useCase) validateInput(id int, status enums.ServiceStatus) errors.Error {
	var err errors.Error

	if errID := uc.validateID(id); errID != nil {
		err = errors.Aggregate(err, errID)
	}

	if errStatus := uc.validateStatus(status); errStatus != nil {
		err = errors.Aggregate(err, errStatus)
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

func (uc *useCase) validateStatus(status enums.ServiceStatus) errors.Error {
	return uc.validator.ValidateVariable(
		status,
		"status",
		"required,enums=enabled disabled deprecated",
		map[string]string{
			"enums":    "status must be one of the following: enabled, disabled, deprecated",
			"required": "status is required",
		},
	)
}

func NewUseCase(
	validator validator.Validator, serviceRepo ServiceRepository,
) UseCase {
	return &useCase{
		serviceRepo: serviceRepo,
		validator:   validator,
	}
}
