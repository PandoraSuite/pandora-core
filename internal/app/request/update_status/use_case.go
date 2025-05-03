package updatestatus

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, id string, executionStatus enums.RequestLogExecutionStatus) errors.Error
}

type useCase struct {
	validator validator.Validator

	requestRepo RequestRepository
}

func (uc *useCase) Execute(
	ctx context.Context,
	id string,
	executionStatus enums.RequestLogExecutionStatus,
) errors.Error {
	if err := uc.validateInput(id, executionStatus); err != nil {
		return err
	}

	return uc.requestRepo.UpdateExecutionStatus(ctx, id, executionStatus)
}

func (uc *useCase) validateInput(
	id string,
	executionStatus enums.RequestLogExecutionStatus,
) errors.Error {
	var err errors.Error

	if errID := uc.validateID(id); errID != nil {
		err = errors.Aggregate(err, errID)
	}

	if errStatus := uc.validateExecutionStatus(executionStatus); errStatus != nil {
		err = errors.Aggregate(err, errStatus)
	}

	return err
}

func (uc *useCase) validateID(id string) errors.Error {
	return uc.validator.ValidateVariable(
		id,
		"id",
		"required,uuid4",
		map[string]string{
			"uuid4":    "id must be a valid UUID",
			"required": "id is required",
		},
	)
}

func (uc *useCase) validateExecutionStatus(executionStatus enums.RequestLogExecutionStatus) errors.Error {
	return uc.validator.ValidateVariable(
		executionStatus,
		"executionStatus",
		"required,enums=success,failed,server_error",
		map[string]string{
			"enums":    "status must be one of the following: success, failed, server_error",
			"required": "executionStatus is required",
		},
	)
}

func NewUseCase(
	validator validator.Validator, requestRepo RequestRepository,
) UseCase {
	return &useCase{
		validator:   validator,
		requestRepo: requestRepo,
	}
}
