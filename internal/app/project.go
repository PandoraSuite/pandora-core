package app

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	domainErr "github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/ports/outbound"
)

type ProjectUseCase struct {
	projectRepo        outbound.ProjectPort
	projectServiceRepo outbound.ProjectServicePort
}

func (u *ProjectUseCase) AssignService(
	ctx context.Context, req *dto.AssignServiceToProject,
) error {
	projectService := &entities.ProjectService{
		ProjectID:      req.ProjectID,
		ServiceID:      req.ServiceID,
		MaxRequest:     req.MaxRequest,
		ResetFrequency: req.ResetFrequency,
	}

	projectService.CalculateNextReset()
	_, err := u.projectServiceRepo.Save(ctx, projectService)
	if err != nil {
		return err
	}

	return nil
}

func (u *ProjectUseCase) GetProjectsByClient(
	ctx context.Context, clientID int,
) ([]*dto.ProjectResponse, error) {
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
) (*dto.ProjectResponse, error) {
	if req.Name == "" {
		return nil, domainErr.ErrNameCannotBeEmpty
	}

	client, err := u.projectRepo.Save(
		ctx,
		&entities.Project{
			Name:     req.Name,
			Status:   req.Status,
			ClientID: req.ClientID,
		},
	)
	if err != nil {
		return nil, err
	}

	return &dto.ProjectResponse{
		ID:        client.ID,
		Name:      client.Name,
		Status:    client.Status,
		ClientID:  client.ClientID,
		CreatedAt: client.CreatedAt,
	}, nil
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
