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
	Execute(ctx context.Context, req *dto.EnvironmentCreate) (*dto.EnvironmentResponse, errors.Error)
}

type useCase struct {
	validator validator.Validator

	projectRepo     ProjectRepository
	environmentRepo EnvironmentRepository
}

func (uc *useCase) Execute(
	ctx context.Context, req *dto.EnvironmentCreate,
) (*dto.EnvironmentResponse, errors.Error) {
	if err := uc.validateReq(req); err != nil {
		return nil, err
	}

	exists, err := uc.projectRepo.Exists(ctx, req.ProjectID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.NewEntityNotFound(
			"Project",
			"project not found",
			map[string]any{"id": req.ProjectID},
			nil,
		)
	}

	var errs errors.Error
	services := make([]*entities.EnvironmentService, len(req.Services))
	for i, service := range req.Services {
		quota, err := uc.projectRepo.GetProjectServiceQuotaUsage(
			ctx, req.ProjectID, service.ID,
		)
		if err != nil {
			if err.Code() == errors.CodeNotFound {
				err = errors.NewEntityNotFound(
					"Service",
					"service not assigned to project",
					map[string]any{"id": service.ID},
					err,
				)
			}
			errs = errors.Aggregate(errs, err)
			continue
		}

		if quota.MaxAllowed > -1 {
			if service.MaxRequests == -1 {
				errs = errors.Aggregate(
					errs,
					errors.NewAttributeValidationFailed(
						"EnvironmentCreate",
						fmt.Sprintf("services[%d].max_requests", i),
						"max_requests cannot be -1 (unlimited) if the project has a defined limit",
						nil,
					),
				)
				continue
			}

			if quota.CurrentAllocated+service.MaxRequests > quota.MaxAllowed {
				errs = errors.Aggregate(
					errs,
					errors.NewAttributeValidationFailed(
						"EnvironmentCreate",
						fmt.Sprintf("services[%d].max_requests", i),
						"max_requests exceeded for service in project",
						nil,
					),
				)
				continue
			}
		}

		services[i] = &entities.EnvironmentService{
			ID:               service.ID,
			MaxRequests:      service.MaxRequests,
			AvailableRequest: service.MaxRequests,
		}
	}

	if errs != nil {
		return nil, errs
	}

	environment := entities.Environment{
		Name:      req.Name,
		Status:    enums.EnvironmentStatusEnabled,
		ProjectID: req.ProjectID,
		Services:  services,
	}

	if err := uc.environmentRepo.Create(ctx, &environment); err != nil {
		return nil, err
	}

	serviceResp := make(
		[]*dto.EnvironmentServiceResponse, len(environment.Services),
	)
	for i, service := range environment.Services {
		serviceResp[i] = &dto.EnvironmentServiceResponse{
			ID:               service.ID,
			Name:             service.Name,
			Version:          service.Version,
			MaxRequests:      service.MaxRequests,
			AvailableRequest: service.AvailableRequest,
			AssignedAt:       service.AssignedAt,
		}
	}

	return &dto.EnvironmentResponse{
		ID:        environment.ID,
		Name:      environment.Name,
		Status:    environment.Status,
		ProjectID: environment.ProjectID,
		CreatedAt: environment.CreatedAt,
		Services:  serviceResp,
	}, nil
}

func (uc *useCase) validateReq(req *dto.EnvironmentCreate) errors.Error {
	return uc.validator.ValidateStruct(
		req,
		map[string]string{
			"name.required":               "name is required",
			"project_id.gt":               "project_id must be greater than 0",
			"services[].id.gt":            "id must be greater than 0",
			"project_id.required":         "project_id is required",
			"services[].id.required":      "id is required",
			"services[].max_requests.gte": "max_requests must be greater than or equal to -1",
		},
	)
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
