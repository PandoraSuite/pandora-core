package listprojects

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, id int) ([]*dto.ProjectResponse, errors.Error)
}

type useCase struct {
	validator validator.Validator

	clientRepo  ClientRepository
	projectRepo ProjectRepository
}

func (uc *useCase) Execute(
	ctx context.Context, id int,
) ([]*dto.ProjectResponse, errors.Error) {
	if err := uc.validateID(id); err != nil {
		return nil, err
	}

	exists, err := uc.clientRepo.Exists(ctx, id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.NewEntityNotFound(
			"client",
			"client not found",
			map[string]any{"id": id},
			nil,
		)
	}

	projects, err := uc.projectRepo.ListByClient(ctx, id)
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

func NewUseCase(
	validator validator.Validator,
	clientRepo ClientRepository,
	projectRepo ProjectRepository,
) UseCase {
	return &useCase{
		validator:   validator,
		clientRepo:  clientRepo,
		projectRepo: projectRepo,
	}
}
