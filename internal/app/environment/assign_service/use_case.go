package assignservice

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, id int, req *dto.EnvironmentService) (*dto.EnvironmentServiceResponse, errors.Error)
}

type useCase struct {
	validator validator.Validator

	environmentRepo EnvironmentRepository
}

func (uc *useCase) Execute(
	ctx context.Context, id int, req *dto.EnvironmentService,
) (*dto.EnvironmentServiceResponse, errors.Error) {
	if err := uc.validateInput(id, req); err != nil {
		return nil, err
	}

	service := entities.EnvironmentService{
		ID:               req.ID,
		MaxRequest:       req.MaxRequest,
		AvailableRequest: req.MaxRequest,
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

	exists, err = uc.environmentRepo.ExistsServiceIn(ctx, id, service.ID)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.NewEntityAlreadyExists(
			"EnvironmentService",
			"service already assigned to environment",
			map[string]any{"id": service.ID},
		)
	}

	quota, err := uc.environmentRepo.GetProjectServiceQuotaUsage(
		ctx, id, service.ID,
	)
	if err != nil {
		if err.Code() == errors.CodeNotFound {
			return nil, errors.NewEntityNotFound(
				"Service",
				"service not assigned to project",
				map[string]any{"id": service.ID},
			)
		}
		return nil, err
	}

	if quota.MaxAllowed > -1 {
		if service.MaxRequest == -1 {
			return nil, errors.NewAttributeValidationFailed(
				"EnvironmentService",
				"max_request",
				"max_request cannot be -1 (unlimited) if the project has a defined limit",
				nil,
			)
		}

		if quota.CurrentAllocated+service.MaxRequest > quota.MaxAllowed {
			return nil, errors.NewAttributeValidationFailed(
				"EnvironmentService",
				"max_request",
				"max_request exceeded for service in project",
				nil,
			)
		}
	}

	if err := uc.environmentRepo.AddService(ctx, id, &service); err != nil {
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

func (uc *useCase) validateInput(id int, req *dto.EnvironmentService) errors.Error {
	var err errors.Error

	if errID := uc.validateID(id); errID != nil {
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

func (uc *useCase) validateReq(req *dto.EnvironmentService) errors.Error {
	return uc.validator.ValidateStruct(
		req,
		map[string]string{
			"id.gt":           "id must be greater than 0",
			"id.required":     "id is required",
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
