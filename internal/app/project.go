package app

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/ports/outbound"
)

type ProjectUseCase struct {
	projectRepo     outbound.ProjectPort
	environmentRepo outbound.EnvironmentPort
}

func (u *ProjectUseCase) AssignService(
	ctx context.Context, id int, req *dto.ProjectService,
) *errors.Error {
	service := &entities.ProjectService{
		ID:             req.ID,
		MaxRequest:     req.MaxRequest,
		ResetFrequency: req.ResetFrequency,
	}

	if err := service.Validate(); err != nil {
		return err
	}

	service.CalculateNextReset()
	return u.projectRepo.AddService(ctx, id, service)
}

func (u *ProjectUseCase) RemoveService(
	ctx context.Context, id, serviceID int,
) *errors.Error {
	exists, err := u.projectRepo.Exists(ctx, id)
	if err != nil {
		return err
	}

	if !exists {
		return errors.ErrProjectNotFound
	}

	_, err = u.environmentRepo.RemoveServiceFromProjectEnvironments(
		ctx, id, serviceID,
	)
	if err != nil {
		return err
	}

	n, err := u.projectRepo.RemoveServiceFromProject(
		ctx, id, serviceID,
	)
	if err != nil {
		return err
	}

	if n == 0 {
		return errors.ErrServiceNotFound
	}

	return nil
}

func (u *ProjectUseCase) GetByID(
	ctx context.Context, id int,
) (*dto.ProjectResponse, *errors.Error) {
	project, err := u.projectRepo.FindByID(ctx, id)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrProjectNotFound
		}
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

func (u *ProjectUseCase) GetEnvironments(
	ctx context.Context, id int,
) ([]*dto.EnvironmentResponse, *errors.Error) {
	exists, err := u.projectRepo.Exists(ctx, id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.ErrProjectNotFound
	}

	environments, err := u.environmentRepo.FindByProject(ctx, id)
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
				ID:               service.ID,
				Name:             service.Name,
				Version:          service.Version,
				MaxRequest:       service.MaxRequest,
				AvailableRequest: service.AvailableRequest,
				AssignedAt:       service.AssignedAt,
			}
		}

		environmentResponses[i] = &dto.EnvironmentResponse{
			ID:        environment.ID,
			Name:      environment.Name,
			Status:    environment.Status,
			ProjectID: environment.ProjectID,
			CreatedAt: environment.CreatedAt,
			Services:  serviceResp,
		}
	}

	return environmentResponses, nil
}

func (u *ProjectUseCase) Create(
	ctx context.Context, req *dto.ProjectCreate,
) (*dto.ProjectResponse, *errors.Error) {
	services := make([]*entities.ProjectService, len(req.Services))
	for i, service := range req.Services {
		services[i] = &entities.ProjectService{
			ID:             service.ID,
			MaxRequest:     service.MaxRequest,
			ResetFrequency: service.ResetFrequency,
		}
	}

	project := entities.Project{
		Name:     req.Name,
		Status:   req.Status,
		ClientID: req.ClientID,
		Services: services,
	}

	if err := project.Validate(); err != nil {
		return nil, err
	}

	project.CalculateNextServicesReset()

	if err := u.projectRepo.Save(ctx, &project); err != nil {
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

func NewProjectUseCase(
	projectRepo outbound.ProjectPort,
	environmentRepo outbound.EnvironmentPort,
) *ProjectUseCase {
	return &ProjectUseCase{
		projectRepo:     projectRepo,
		environmentRepo: environmentRepo,
	}
}
