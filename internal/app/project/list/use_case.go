package list

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type UseCase interface {
	Execute(ctx context.Context) ([]*dto.ProjectResponse, errors.Error)
}

type useCase struct {
	projectRepo ProjectRepository
}

func (uc *useCase) Execute(ctx context.Context) ([]*dto.ProjectResponse, errors.Error) {
	projects, err := uc.projectRepo.List(ctx)
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
				Name:           service.Name,
				Version:        service.Version,
				NextReset:      service.NextReset,
				MaxRequests:    service.MaxRequests,
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

func NewUseCase(projectRepo ProjectRepository) UseCase {
	return &useCase{
		projectRepo: projectRepo,
	}
}
