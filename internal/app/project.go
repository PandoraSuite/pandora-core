package app

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/ports/outbound"
)

type ProjectUseCase struct {
	projectRepo        outbound.ProjectPort
	projectServiceRepo outbound.ProjectServicePort
}

func (u *ProjectUseCase) AssignService(
	ctx context.Context, req *dto.AssignServiceToProject,
) *errors.Error {
	projectService := &entities.ProjectService{
		ProjectID:      req.ProjectID,
		ServiceID:      req.ServiceID,
		MaxRequest:     req.MaxRequest,
		ResetFrequency: req.ResetFrequency,
	}

	projectService.CalculateNextReset()
	if err := u.projectServiceRepo.Save(ctx, projectService); err != nil {
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
		projectResponses[i] = &dto.ProjectResponse{
			ID:        project.ID,
			Name:      project.Name,
			Status:    project.Status,
			ClientID:  project.ClientID,
			CreatedAt: project.CreatedAt,
		}
	}

	return projectResponses, nil
}

func (u *ProjectUseCase) Create(
	ctx context.Context, req *dto.ProjectCreate,
) (*dto.ProjectResponse, *errors.Error) {
	project := entities.Project{
		Name:     req.Name,
		Status:   req.Status,
		ClientID: req.ClientID,
	}

	if err := project.Validate(); err != nil {
		return nil, err
	}

	if err := u.projectRepo.Save(ctx, &project); err != nil {
		return nil, err
	}

	resp := &dto.ProjectResponse{
		ID:        project.ID,
		Name:      project.Name,
		Status:    project.Status,
		ClientID:  project.ClientID,
		CreatedAt: project.CreatedAt,
	}

	var servicesErrors []string
	var projectServices []*entities.ProjectService
	for _, s := range req.Services {
		projectService := &entities.ProjectService{
			ProjectID:      project.ID,
			ServiceID:      s.ID,
			MaxRequest:     s.MaxRequest,
			ResetFrequency: s.ResetFrequency,
		}

		if err := projectService.Validate(); err != nil {
			servicesErrors = append(servicesErrors, err.Error())
			continue
		}

		projectService.CalculateNextReset()
		projectServices = append(projectServices, projectService)
	}

	var errResp *errors.Error
	if len(servicesErrors) > 0 {
		errResp = errors.NewError(
			errors.CodeValidationError,
			"Invalid service assignments",
			servicesErrors...,
		)
	}

	if len(projectServices) > 0 {
		err := u.projectServiceRepo.BulkSave(ctx, projectServices)
		if err != nil {
			if errResp != nil {
				err = errors.NewError(
					err.Code,
					"Failed to save project services",
					append(servicesErrors, err.Message)...,
				)
			}
			return resp, err
		}

		projectServiceResp := make(
			[]*dto.ProjectServiceAssignmentResponse,
			len(projectServices),
		)
		for i, projectService := range projectServices {
			projectServiceResp[i] = &dto.ProjectServiceAssignmentResponse{
				ID:             projectService.ServiceID,
				NextReset:      projectService.NextReset,
				MaxRequest:     projectService.MaxRequest,
				ResetFrequency: projectService.ResetFrequency,
			}
		}
		resp.Services = projectServiceResp
	}

	return resp, errResp
}

func NewProjectUseCase(
	projectRepo outbound.ProjectPort,
	projectServiceRepo outbound.ProjectServicePort,
) *ProjectUseCase {
	return &ProjectUseCase{
		projectRepo:        projectRepo,
		projectServiceRepo: projectServiceRepo,
	}
}
