package get

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, id int) (*dto.EnvironmentResponse, errors.Error)
}

type useCase struct {
	validator validator.Validator

	environmentRepo EnvironmentRepository
}

func (uc *useCase) Execute(ctx context.Context, id int) (*dto.EnvironmentResponse, errors.Error) {
	if err := uc.validateID(id); err != nil {
		return nil, err
	}

	environment, err := uc.environmentRepo.GetByID(ctx, id)
	if err != nil {
		if err.Code() == errors.CodeNotFound {
			return nil, errors.NewEntityNotFound(
				"Environment",
				"environment not found",
				map[string]any{"id": id},
				err,
			)
		}
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
			MaxRequest:       service.MaxRequest,
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
	validator validator.Validator, environmentRepo EnvironmentRepository,
) UseCase {
	return &useCase{
		validator:       validator,
		environmentRepo: environmentRepo,
	}
}
