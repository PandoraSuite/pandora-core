package app

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/ports/outbound"
)

type EnvironmentUseCase struct {
	environmentRepo outbound.EnvironmentPort
}

func (u *EnvironmentUseCase) AssignService(
	ctx context.Context, id int, req *dto.EnvironmentService,
) *errors.Error {
	service := entities.EnvironmentService{
		ID:               req.ID,
		MaxRequest:       req.MaxRequest,
		AvailableRequest: req.MaxRequest,
	}

	if id <= 0 {
		return errors.ErrInvalidEnvironmentID
	}

	if err := service.Validate(); err != nil {
		return err
	}

	exists, err := u.environmentRepo.ExistsEnvironmentService(ctx, id, service.ID)
	if err != nil {
		return err
	}

	if exists {
		return errors.ErrEnvironmentServiceAlreadyExists
	}

	maxRequest, err := u.environmentRepo.
		GetMaxRequestForServiceInProject(ctx, id, service.ID)
	if err != nil {
		return err
	}

	if maxRequest > 0 {
		maxRequests, err := u.environmentRepo.
			GetAllMaxRequestForServiceInEnvironments(ctx, id, req.ID)
		if err != nil {
			return err
		}

		var totalMaxRequest int
		for _, v := range maxRequests {
			totalMaxRequest += v
		}

		if req.MaxRequest+totalMaxRequest > maxRequest {
			return errors.ErrMaxRequestExceededForServiceInProyect
		}
	}

	return u.environmentRepo.AddService(ctx, id, &service)
}

func (u *EnvironmentUseCase) GetByProject(
	ctx context.Context, projectID int,
) ([]*dto.EnvironmentResponse, *errors.Error) {
	environments, err := u.environmentRepo.FindByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	environmentResponses := make([]*dto.EnvironmentResponse, len(environments))
	for i, environment := range environments {
		serviceResp := make(
			[]*dto.EnvironmentServiceResponse, len(environment.Services),
		)
		for i, service := range environment.Services {
			serviceResp[i] = &dto.EnvironmentServiceResponse{
				ID:         service.ID,
				Name:       service.Name,
				Version:    service.Version,
				MaxRequest: service.MaxRequest,
				AssignedAt: service.AssignedAt,
			}
		}

		environmentResponses[i] = &dto.EnvironmentResponse{
			ID:        environment.ID,
			Name:      environment.Name,
			Status:    environment.Status,
			ProjectID: environment.ProjectID,
			CreatedAt: environment.CreatedAt,
		}
	}

	return environmentResponses, nil
}

func (u *EnvironmentUseCase) Create(
	ctx context.Context, req *dto.EnvironmentCreate,
) (*dto.EnvironmentResponse, *errors.Error) {
	services := make([]*entities.EnvironmentService, len(req.Services))
	for i, service := range req.Services {
		services[i] = &entities.EnvironmentService{
			ID:               service.ID,
			MaxRequest:       service.MaxRequest,
			AvailableRequest: service.MaxRequest,
		}
	}

	environment := entities.Environment{
		Name:      req.Name,
		Status:    enums.EnvironmentActive,
		ProjectID: req.ProjectID,
		Services:  services,
	}

	if err := environment.Validate(); err != nil {
		return nil, err
	}

	if err := u.environmentRepo.Save(ctx, &environment); err != nil {
		return nil, err
	}

	serviceResp := make(
		[]*dto.EnvironmentServiceResponse, len(environment.Services),
	)
	for i, service := range environment.Services {
		serviceResp[i] = &dto.EnvironmentServiceResponse{
			ID:         service.ID,
			Name:       service.Name,
			Version:    service.Version,
			MaxRequest: service.MaxRequest,
			AssignedAt: service.AssignedAt,
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

func NewEnvironmentUseCase(environmentRepo outbound.EnvironmentPort) *EnvironmentUseCase {
	return &EnvironmentUseCase{environmentRepo: environmentRepo}
}
