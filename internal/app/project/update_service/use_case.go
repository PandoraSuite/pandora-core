package updateservice

import (
	"context"
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/utils"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, id, serviceID int, req *dto.ProjectServiceUpdate) (*dto.ProjectServiceResponse, errors.Error)
}

type useCase struct {
	validator validator.Validator

	projectRepo     ProjectRepository
	environmentRepo EnvironmentRepository
}

func (uc *useCase) Execute(
	ctx context.Context, id, serviceID int, req *dto.ProjectServiceUpdate,
) (*dto.ProjectServiceResponse, errors.Error) {

	if req.NextReset.IsZero() {
		service := entities.ProjectService{
			MaxRequest:     req.MaxRequest,
			ResetFrequency: req.ResetFrequency,
		}

		service.CalculateNextReset()
		req.NextReset = service.NextReset
	} else {
		req.NextReset = utils.TruncateToDay(req.NextReset)
	}

	if err := uc.validateInput(id, serviceID, req); err != nil {
		return nil, err
	}

	exists, err := uc.projectRepo.Exists(ctx, id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.NewEntityNotFound(
			"Project",
			"project not found",
			map[string]any{"id": id},
			nil,
		)
	}

	quota, err := uc.projectRepo.GetProjectServiceQuotaUsage(
		ctx, id, serviceID,
	)
	if err != nil {
		if err.Code() == errors.CodeNotFound {
			return nil, errors.NewEntityNotFound(
				"Service",
				"service not assigned to project",
				map[string]any{"id": serviceID},
				err,
			)
		}
		return nil, err
	}

	if req.MaxRequest != -1 {
		if req.MaxRequest < quota.CurrentAllocated {
			return nil, errors.NewAttributeValidationFailed(
				"ProjectServiceUpdate",
				"max_request",
				"max_request is below the total allocated to environments",
				nil,
			)
		}
		hasInfinite, err := uc.environmentRepo.ExistsServiceWithInfiniteMaxRequest(
			ctx, id, serviceID,
		)
		if err != nil {
			return nil, err
		}

		if hasInfinite {
			return nil, errors.NewValidationFailed(
				"cannot set a finite max_request while some environments have infinite quota",
				nil,
			)
		}
	}

	service, err := uc.projectRepo.UpdateService(ctx, id, serviceID, req)
	if err != nil {
		return nil, err
	}

	return &dto.ProjectServiceResponse{
		ID:             service.ID,
		Name:           service.Name,
		Version:        service.Version,
		NextReset:      service.NextReset,
		MaxRequest:     service.MaxRequest,
		ResetFrequency: service.ResetFrequency,
		AssignedAt:     service.AssignedAt,
	}, nil
}

func (uc *useCase) validateInput(id, serviceID int, req *dto.ProjectServiceUpdate) errors.Error {
	var err errors.Error

	if errID := uc.validateID(id); errID != nil {
		err = errors.Aggregate(err, errID)
	}

	if errServiceID := uc.validateServiceID(serviceID); errServiceID != nil {
		err = errors.Aggregate(err, errServiceID)
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

func (uc *useCase) validateReq(req *dto.ProjectServiceUpdate) errors.Error {
	var err errors.Error

	validationErr := uc.validator.ValidateStruct(
		req,
		map[string]string{
			"next_reset.utc":        "next_reset must be a valid UTC datetime",
			"max_request.gte":       "max_request must be greater than or equal to -1",
			"reset_frequency.enums": "reset_frequency must be one of the following: , daily, weekly, biweekly, monthly",
		},
	)

	if validationErr != nil {
		err = errors.Aggregate(err, validationErr)
	}

	if req.MaxRequest == -1 && req.ResetFrequency != enums.ProjectServiceNull {
		err = errors.Aggregate(
			err,
			errors.NewAttributeValidationFailed(
				"ProjectServiceUpdate",
				"reset_frequency",
				"reset_frequency must be null when max_request is -1 (unlimited)",
				nil,
			),
		)
	}

	if req.MaxRequest > -1 && req.ResetFrequency == enums.ProjectServiceNull {
		err = errors.Aggregate(
			err,
			errors.NewAttributeValidationFailed(
				"ProjectServiceUpdate",
				"reset_frequency",
				"reset_frequency is required when max_request is greater than -1 (unlimited)",
				nil,
			),
		)
	}

	if !req.NextReset.IsZero() {
		if req.NextReset.Before(utils.TruncateToDay(time.Now())) {
			err = errors.Aggregate(
				err,
				errors.NewAttributeValidationFailed(
					"ProjectServiceUpdate",
					"next_reset",
					"next_reset must be in the future",
					nil,
				),
			)
		}
		if req.MaxRequest == -1 {
			err = errors.Aggregate(
				err,
				errors.NewAttributeValidationFailed(
					"ProjectServiceUpdate",
					"next_reset",
					"next_reset must be null when max_request is -1 (unlimited)",
					nil,
				),
			)
		}
	}

	return err
}

func NewUseCase(
	validator validator.Validator,
	projectRepo ProjectRepository,
	environmentRepo EnvironmentRepository,
) UseCase {
	return &useCase{
		validator:       validator,
		projectRepo:     projectRepo,
		environmentRepo: environmentRepo,
	}
}
