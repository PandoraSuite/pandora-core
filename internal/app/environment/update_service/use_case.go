package updateservice

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, id, serviceID int, req *dto.EnvironmentServiceUpdateInput) (*dto.EnvironmentServiceResponse, errors.Error)
}

type useCase struct {
	validator validator.Validator

	environmentRepo EnvironmentRepository
}

func (uc *useCase) Execute(
	ctx context.Context, id, serviceID int, req *dto.EnvironmentServiceUpdateInput,
) (*dto.EnvironmentServiceResponse, errors.Error) {
	if err := uc.validateInput(id, serviceID, req); err != nil {
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
		)
	}

	service, err := uc.environmentRepo.GetServiceByID(
		ctx, id, serviceID,
	)
	if err != nil {
		if err.Code() == errors.CodeNotFound {
			return nil, errors.NewEntityNotFound(
				"Service",
				"service not assigned to environment",
				map[string]any{"id": serviceID},
			)
		}
		return nil, err
	}

	quota, err := uc.environmentRepo.GetProjectServiceQuotaUsage(
		ctx, id, serviceID,
	)
	if err != nil {
		if err.Code() == errors.CodeNotFound {
			return nil, errors.NewEntityNotFound(
				"Service",
				"service not assigned to project",
				map[string]any{"id": serviceID},
			)
		}
		return nil, err
	}

	if quota.MaxAllowed != -1 {
		if req.MaxRequest == -1 {
			return nil, errors.NewAttributeValidationFailed(
				"EnvironmentServiceUpdateInput",
				"max_request",
				"max_request cannot be -1 (unlimited) if the project has a defined limit",
				nil,
			)
		}

		if quota.CurrentAllocated-service.MaxRequest+req.MaxRequest > quota.MaxAllowed {
			return nil, errors.NewAttributeValidationFailed(
				"EnvironmentCreate",
				"max_request",
				"max_request exceeded for service in project",
				nil,
			)
		}
	}

	update := &dto.EnvironmentServiceUpdate{MaxRequest: req.MaxRequest}
	if service.AvailableRequest == -1 || service.AvailableRequest > req.MaxRequest {
		update.AvailableRequest = req.MaxRequest
	} else {
		update.AvailableRequest = service.AvailableRequest
	}

	service, err = uc.environmentRepo.UpdateService(
		ctx, id, serviceID, update,
	)
	if err != nil {
		return nil, err
	}

	return &dto.EnvironmentServiceResponse{
		ID:               service.ID,
		Name:             service.Name,
		Version:          service.Version,
		MaxRequest:       service.MaxRequest,
		AvailableRequest: service.AvailableRequest,
		AssignedAt:       service.AssignedAt,
	}, nil
}

func (uc *useCase) validateInput(
	id, serviceID int, req *dto.EnvironmentServiceUpdateInput,
) errors.Error {
	var err errors.Error

	if errID := uc.validateID(id); errID != nil {
		err = errors.Aggregate(err, errID)
	}

	if errID := uc.validateServiceID(serviceID); errID != nil {
		err = errors.Aggregate(err, errID)
	}

	if errReq := uc.validateReq(req); errReq != nil {
		err = errors.Aggregate(err, errReq)
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

func (uc *useCase) validateReq(req *dto.EnvironmentServiceUpdateInput) errors.Error {
	return uc.validator.ValidateStruct(
		req,
		map[string]string{
			"max_request.gte": "max_request must be greater than or equal to -1",
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
