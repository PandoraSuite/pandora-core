package listrequest

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, id int, req *dto.RequestFilter) ([]*dto.RequestResponse, errors.Error)
}

type useCase struct {
	validator validator.Validator

	serviceRepo ServiceRepository
	requestRepo RequestRepository
}

func (uc *useCase) Execute(
	ctx context.Context, id int, req *dto.RequestFilter,
) ([]*dto.RequestResponse, errors.Error) {
	if err := uc.validateInput(id, req); err != nil {
		return nil, err
	}

	exists, err := uc.serviceRepo.Exists(ctx, id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.NewEntityNotFound(
			"Service",
			"service not found",
			map[string]any{"id": id},
			nil,
		)
	}

	requests, err := uc.requestRepo.ListByService(ctx, id, req)
	if err != nil {
		return nil, err
	}

	requestResponses := make([]*dto.RequestResponse, len(requests))
	for i, request := range requests {
		requestResponses[i] = &dto.RequestResponse{
			ID:                 request.ID,
			StartPoint:         request.StartPoint,
			StatusCode:         request.StatusCode,
			ExecutionStatus:    request.ExecutionStatus,
			UnauthorizedReason: request.UnauthorizedReason,
			RequestTime:        req.RequestTimeFrom,
			Path:               request.Path,
			Method:             request.Method,
			IPAddress:          request.IPAddress,
			CreatedAt:          request.CreatedAt,
			APIKey: &dto.RequestAPIKeyResponse{
				ID:  request.APIKey.ID,
				Key: request.APIKey.KeySummary(),
			},
			Project: &dto.RequestProjectResponse{
				ID:   request.Project.ID,
				Name: request.Project.Name,
			},
			Environment: &dto.RequestEnvironmentResponse{
				ID:   request.Environment.ID,
				Name: request.Environment.Name,
			},
			Service: &dto.RequestServiceResponse{
				ID:      request.Service.ID,
				Name:    request.Service.Name,
				Version: request.Service.Version,
			},
		}
	}

	return requestResponses, nil
}

func (uc *useCase) validateInput(id int, req *dto.RequestFilter) errors.Error {
	var err errors.Error

	if errID := uc.validateID(id); errID != nil {
		err = errors.Aggregate(err, errID)
	}

	if errStatus := uc.validateReq(req); errStatus != nil {
		err = errors.Aggregate(err, errStatus)
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

func (uc *useCase) validateReq(req *dto.RequestFilter) errors.Error {
	return uc.validator.ValidateStruct(
		req,
		map[string]string{
			"request_time_from.gtetime": "request_time_from must be greater than request_time_to",
			"request_time_to":           "request_time_to is required",
			"execution_status.enums":    "execution_status must be one of the following: success, forwarded, client_error, server_error, unauthorized, quota_exceeded",
		},
	)
}

func NewUseCase(
	validator validator.Validator,
	serviceRepo ServiceRepository,
	requestRepo RequestRepository,
) UseCase {
	return &useCase{
		validator:   validator,
		serviceRepo: serviceRepo,
		requestRepo: requestRepo,
	}
}
