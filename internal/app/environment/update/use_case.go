package update

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, id int, req *dto.EnvironmentUpdate) (*dto.EnvironmentResponse, errors.Error)
}

type useCase struct {
	validator validator.Validator

	environmentRepo EnvironmentRepository
}

func (uc *useCase) Execute(
	ctx context.Context, id int, req *dto.EnvironmentUpdate,
) (*dto.EnvironmentResponse, errors.Error) {
	if err := uc.validatorInput(id, req); err != nil {
		return nil, err
	}

	environment, err := uc.environmentRepo.Update(ctx, id, req)
	if err != nil {
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
			MaxRequests:      service.MaxRequests,
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

func (uc *useCase) validatorInput(id int, req *dto.EnvironmentUpdate) errors.Error {
	var err errors.Error

	if errID := uc.validateID(id); errID != nil {
		err = errors.Aggregate(err, errID)
	}

	if errReq := uc.validateReq(req); errReq != nil {
		err = errors.Aggregate(err, errReq)
	}

	return err
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

func (uc *useCase) validateReq(req *dto.EnvironmentUpdate) errors.Error {
	return uc.validator.ValidateStruct(req, map[string]string{})
}

func NewUseCase(
	validator validator.Validator, environmentRepo EnvironmentRepository,
) UseCase {
	return &useCase{
		validator:       validator,
		environmentRepo: environmentRepo,
	}
}
