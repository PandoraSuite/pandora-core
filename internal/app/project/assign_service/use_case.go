package assignservice

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, id int, req *dto.ProjectService) (*dto.ProjectServiceResponse, errors.Error)
}

type useCase struct {
	validator validator.Validator

	projectRepo ProjectRepository
}

func (uc *useCase) Execute(
	ctx context.Context, id int, req *dto.ProjectService,
) (*dto.ProjectServiceResponse, errors.Error) {
	if err := uc.validateInput(id, req); err != nil {
		return nil, err
	}

	service := &entities.ProjectService{
		ID:             req.ID,
		MaxRequest:     req.MaxRequest,
		ResetFrequency: req.ResetFrequency,
	}

	service.CalculateNextReset()
	if err := uc.projectRepo.AddService(ctx, id, service); err != nil {
		if err.Code() == errors.CodeAlreadyExists {
			return nil, errors.NewEntityAlreadyExists(
				"ProjectService",
				"service already assigned to project",
				map[string]any{"id": service.ID},
				err,
			)
		}
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

func (uc *useCase) validateInput(id int, req *dto.ProjectService) errors.Error {
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

func (uc *useCase) validateReq(req *dto.ProjectService) errors.Error {
	var err errors.Error

	validationErr := uc.validator.ValidateStruct(
		req,
		map[string]string{
			"id.gt":                 "id must be greater than 0",
			"id.required":           "id is required",
			"max_request.gte":       "max_request must be greater than or equal to -1",
			"reset_frequency.enums": "reset_frequency must be one of the following: , daily, weekly, biweekly, monthly",
		},
	)

	if validationErr != nil {
		err = errors.Aggregate(err, validationErr)
	}

	if req.MaxRequest == -1 && req.ResetFrequency != enums.ProjectServiceResetFrequencyNull {
		err = errors.Aggregate(
			err,
			errors.NewAttributeValidationFailed(
				"ProjectService",
				"reset_frequency",
				"reset_frequency must be null when max_request is -1 (unlimited)",
				nil,
			),
		)
	}

	if req.MaxRequest > -1 && req.ResetFrequency == enums.ProjectServiceResetFrequencyNull {
		err = errors.Aggregate(
			err,
			errors.NewAttributeValidationFailed(
				"ProjectService",
				"reset_frequency",
				"reset_frequency is required when max_request is greater than -1 (unlimited)",
				nil,
			),
		)
	}

	return err
}

func NewUseCase(
	validator validator.Validator, projectRepo ProjectRepository,
) UseCase {
	return &useCase{
		validator:   validator,
		projectRepo: projectRepo,
	}
}
