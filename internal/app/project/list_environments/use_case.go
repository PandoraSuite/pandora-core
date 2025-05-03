package listenvironments

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, id int) ([]*dto.EnvironmentResponse, errors.Error)
}

type useCase struct {
	validator validator.Validator

	projectRepo     ProjectRepository
	environmentRepo EnvironmentRepository
}

func (uc *useCase) Execute(
	ctx context.Context, id int,
) ([]*dto.EnvironmentResponse, errors.Error) {
	if err := uc.validateID(id); err != nil {
		return nil, err
	}

	exists, err := uc.projectRepo.Exists(ctx, id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.NewEntityNotFound(
			"Project",
			"project not found",
			map[string]any{"id": id},
		)
	}

	environments, err := uc.environmentRepo.ListByProject(ctx, id)
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
	projectRepo ProjectRepository,
	environmentRepo EnvironmentRepository,
) UseCase {
	return &useCase{
		validator:       validator,
		projectRepo:     projectRepo,
		environmentRepo: environmentRepo,
	}
}
