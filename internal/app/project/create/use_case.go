package create

import (
	"context"

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
	if err := uc.validateReq(req); err != nil {
		return nil, err
	}

	services := make([]*entities.ProjectService, len(req.Services))
	for i, service := range req.Services {
		s := &entities.ProjectService{
			ID:             service.ID,
			MaxRequests:    service.MaxRequests,
			ResetFrequency: service.ResetFrequency,
		}
		s.CalculateNextReset()
		services[i] = s
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
			MaxRequests:    service.MaxRequests,
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
			"name.required":                       "name is required",
			"client_id.gt":                        "client_id must be greater than 0",
			"status.required":                     "status is required",
			"client_id.required":                  "client_id is required",
			"services[].id.gt":                    "id must be greater than 0",
			"services[].id.required":              "id is required",
			"services[].max_requests.gte":         "max_requests must be greater than or equal to -1",
			"services[].reset_frequency.enums":    "reset_frequency must be one of the following: daily, weekly, biweekly, monthly",
			"services[].reset_frequency.required": "reset_frequency is required",
		},
	)
}

func NewUseCase(
	validator validator.Validator, projectRepo ProjectRepository,
) UseCase {
	return &useCase{
		validator:   validator,
		projectRepo: projectRepo,
	}
}
