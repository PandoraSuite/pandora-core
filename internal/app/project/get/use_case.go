package get

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, id int) (*dto.ProjectResponse, errors.Error)
}

type useCase struct {
	validator validator.Validator

	projectRepo ProjectRepository
}

func (uc *useCase) Execute(
	ctx context.Context, id int,
) (*dto.ProjectResponse, errors.Error) {
	if err := uc.validateID(id); err != nil {
		return nil, err
	}

	project, err := uc.projectRepo.GetByID(ctx, id)
	if err != nil {
		if err.Code() == errors.CodeNotFound {
			return nil, errors.NewEntityNotFound(
				"Project",
				"project not found",
				map[string]any{"id": id},
				err,
			)
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
	validator validator.Validator, projectRepo ProjectRepository,
) UseCase {
	return &useCase{
		validator:   validator,
		projectRepo: projectRepo,
	}
}
