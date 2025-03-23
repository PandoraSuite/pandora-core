package app

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/ports/outbound"
)

type ProjectUseCase struct {
	projectRepo outbound.ProjectPort
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
	if err := u.projectRepo.AddService(ctx, id, service); err != nil {
		return err
	}

	return nil
}

func (u *ProjectUseCase) GetByClient(
	ctx context.Context, clientID int,
) ([]*dto.ProjectResponse, *errors.Error) {
	projects, err := u.projectRepo.FindByClient(ctx, clientID)
	if err != nil {
		return nil, err
	}

	projectResponses := make([]*dto.ProjectResponse, len(projects))
	for i, project := range projects {
		serviceResp := make(
			[]*dto.ProjectServiceResponse, len(project.Services),
		)
		for i, service := range project.Services {
			serviceResp[i] = &dto.ProjectServiceResponse{
				ID:             service.ID,
				NextReset:      service.NextReset,
				MaxRequest:     service.MaxRequest,
				ResetFrequency: service.ResetFrequency,
				AssignedAt:     service.AssignedAt,
			}
		}

		projectResponses[i] = &dto.ProjectResponse{
			ID:        project.ID,
			Name:      project.Name,
			Status:    project.Status,
			ClientID:  project.ClientID,
			CreatedAt: project.CreatedAt,
			Services:  serviceResp,
		}
	}

	return projectResponses, nil
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

func NewProjectUseCase(projectRepo outbound.ProjectPort) *ProjectUseCase {
	return &ProjectUseCase{projectRepo: projectRepo}
}
