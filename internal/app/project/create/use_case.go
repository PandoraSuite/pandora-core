package create

import (
	"context"
	"fmt"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, req *dto.ProjectCreate) (*dto.ProjectResponse, errors.Error)
}

type useCase struct {
	validator validator.Validator

	projectRepo ProjectRepository
}

func (uc *useCase) Execute(
	ctx context.Context, req *dto.ProjectCreate,
) (*dto.ProjectResponse, errors.Error) {
	validationErr := uc.validateReq(req)
	services := make([]*entities.ProjectService, len(req.Services))
	for i, service := range req.Services {
		if err := uc.validateService(service); err != nil {
			if err, ok := err.(*errors.AttributeError); ok {
				err.PrefixLoc(fmt.Sprintf("services[%d]", i))
			}
			validationErr = errors.Aggregate(validationErr, err)
			continue
		}

		s := &entities.ProjectService{
			ID:             service.ID,
			MaxRequest:     service.MaxRequest,
			ResetFrequency: service.ResetFrequency,
		}
		s.CalculateNextReset()
		services[i] = s
	}

	if validationErr != nil {
		return nil, validationErr
	}

	project := entities.Project{
		Name:     req.Name,
		Status:   enums.ProjectStatusEnabled,
		ClientID: req.ClientID,
		Services: services,
	}

	if err := uc.projectRepo.Create(ctx, &project); err != nil {
		return nil, err
	}

	serviceResp := make(
		[]*dto.ProjectServiceResponse, len(project.Services),
	)
	for i, service := range project.Services {
		serviceResp[i] = &dto.ProjectServiceResponse{
			ID:             service.ID,
			Name:           service.Name,
			Version:        service.Version,
			NextReset:      service.NextReset,
			MaxRequest:     service.MaxRequest,
			ResetFrequency: service.ResetFrequency,
			AssignedAt:     service.AssignedAt,
		}
	}

	return &dto.ProjectResponse{
		ID:        project.ID,
		Name:      project.Name,
		Status:    project.Status,
		ClientID:  project.ClientID,
		CreatedAt: project.CreatedAt,
		Services:  serviceResp,
	}, nil
}

func (uc *useCase) validateReq(req *dto.ProjectCreate) errors.Error {
	return uc.validator.ValidateStruct(
		req,
		map[string]string{
			"name.required":                    "name is required",
			"client_id.gt":                     "client_id must be greater than 0",
			"status.required":                  "status is required",
			"client_id.required":               "client_id is required",
			"services[].id.gt":                 "id must be greater than 0",
			"services[].id.required":           "id is required",
			"services[].max_request.gte":       "max_request must be greater than or equal to -1",
			"services[].reset_frequency.enums": "reset_frequency must be one of the following: , daily, weekly, biweekly, monthly",
		},
	)
}

func (uc *useCase) validateService(req *dto.ProjectService) errors.Error {
	var err errors.Error

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
