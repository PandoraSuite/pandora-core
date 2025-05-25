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
			"Project",
			"project not found",
			map[string]any{"id": id},
			nil,
		)
	}

	requests, err := uc.requestRepo.ListByService(ctx, id, req)
	if err != nil {
		return nil, err
	}

	return requests, nil
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
			"execution_status.enums":    "execution_status must be one of the following: success, forwarded, client_error, service_error, unauthorized, quota_exceeded",
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
