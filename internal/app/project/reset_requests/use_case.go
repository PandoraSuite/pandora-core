package resetrequests

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, id, serviceID int, req *dto.ProjectServiceResetRequest) (*dto.ProjectServiceResetRequestResponse, errors.Error)
}

type useCase struct {
	validator validator.Validator

	projectRepo ProjectRepository
}

func (uc *useCase) Execute(
	ctx context.Context, id, serviceID int, req *dto.ProjectServiceResetRequest,
) (*dto.ProjectServiceResetRequestResponse, errors.Error) {
	if err := uc.validateInput(id, serviceID, req); err != nil {
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

	service, err := uc.projectRepo.GetServiceByID(ctx, id, serviceID)
	if err != nil {
		if err.Code() == errors.CodeNotFound {
			return nil, errors.NewEntityNotFound(
				"ProjectService",
				"service not assigned to Project",
				map[string]any{"id": serviceID},
			)
		}
		return nil, err
	}

	var envServices []*dto.EnvironmentServiceReset
	if req.RecalculateNextReset {
		service.CalculateNextReset()
		envServices, err = uc.projectRepo.
			ResetProjectServiceUsage(
				ctx, id, serviceID, service.NextReset,
			)
	} else {
		envServices, err = uc.projectRepo.
			ResetAvailableRequestsForEnvsService(
				ctx, id, serviceID,
			)
	}

	if err != nil {
		return nil, err
	}

	return &dto.ProjectServiceResetRequestResponse{
		ResetCount:          len(envServices),
		EnvironmentServices: envServices,
		ProjectService: &dto.ProjectServiceResponse{
			ID:             service.ID,
			Name:           service.Name,
			Version:        service.Version,
			NextReset:      service.NextReset,
			MaxRequest:     service.MaxRequest,
			ResetFrequency: service.ResetFrequency,
			AssignedAt:     service.AssignedAt,
		},
	}, nil
}

func (uc *useCase) validateInput(
	id, serviceID int, req *dto.ProjectServiceResetRequest,
) errors.Error {
	var err errors.Error

	if errID := uc.validateID(id); errID != nil {
		err = errors.Aggregate(err, errID)
	}

	if errServiceID := uc.validateServiceID(serviceID); errServiceID != nil {
		err = errors.Aggregate(err, errServiceID)
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

func (uc *useCase) validateServiceID(serviceID int) errors.Error {
	return uc.validator.ValidateVariable(
		serviceID,
		"service_id",
		"required,gt=0",
		map[string]string{
			"gt":       "service_id must be greater than 0",
			"required": "service_id is required",
		},
	)
}

func (uc *useCase) validateReq(req *dto.ProjectServiceResetRequest) errors.Error {
	return uc.validator.ValidateStruct(req, map[string]string{})
}

func NewUseCase(
	validator validator.Validator, projectRepo ProjectRepository,
) UseCase {
	return &useCase{
		validator:   validator,
		projectRepo: projectRepo,
	}
}
