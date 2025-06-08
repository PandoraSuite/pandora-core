package updateexecutionstatus

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, id string, req *dto.RequestExecutionStatusUpdate) errors.Error
}

type useCase struct {
	validator validator.Validator

	requestRepo RequestRepository
}

func (uc *useCase) Execute(
	ctx context.Context, id string, req *dto.RequestExecutionStatusUpdate,
) errors.Error {
	if err := uc.validateInput(id, req); err != nil {
		return err
	}

	return uc.requestRepo.UpdateExecutionStatus(ctx, id, req)
}

func (uc *useCase) validateInput(
	id string, req *dto.RequestExecutionStatusUpdate,
) errors.Error {
	var err errors.Error

	if errID := uc.validateID(id); errID != nil {
		err = errors.Aggregate(err, errID)
	}

	if errStatus := uc.validateReq(req); errStatus != nil {
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

func (uc *useCase) validateReq(req *dto.RequestExecutionStatusUpdate) errors.Error {
	var err errors.Error

	errReq := uc.validator.ValidateStruct(
		req,
		map[string]string{
			"status_code.required":      "status_code is required",
			"detail.required_unless":    "detail is required unless execution_status is 'success'",
			"execution_status.enums":    "execution_status must be one of the following: success, client_error, server_error",
			"execution_status.required": "execution_status is required",
		},
	)

	if errReq != nil {
		err = errors.Aggregate(err, errReq)
	}

	switch req.ExecutionStatus {
	case enums.RequestExecutionStatusSuccess:
		if req.StatusCode < 200 || req.StatusCode >= 300 {
			errors.Aggregate(
				err,
				errors.NewAttributeValidationFailed(
					"RequestExecutionStatusUpdate",
					"status_code",
					"status_code must be in the range of 200-299 when execution_status is 'success'",
					nil,
				),
			)
		}
	case enums.RequestExecutionStatusClientError:
		if req.StatusCode < 400 || req.StatusCode >= 500 {
			errors.Aggregate(
				err,
				errors.NewAttributeValidationFailed(
					"RequestExecutionStatusUpdate",
					"status_code",
					"status_code must be in the range of 400-499 when execution_status is 'client_error'",
					nil,
				),
			)
		}
	case enums.RequestExecutionStatusServerError:
		if req.StatusCode < 500 || req.StatusCode >= 600 {
			errors.Aggregate(
				err,
				errors.NewAttributeValidationFailed(
					"RequestExecutionStatusUpdate",
					"status_code",
					"status_code must be in the range of 500-599 when execution_status is 'server_error'",
					nil,
				),
			)
		}
	}

	return err
}

func NewUseCase(
	validator validator.Validator, requestRepo RequestRepository,
) UseCase {
	return &useCase{
		validator:   validator,
		requestRepo: requestRepo,
	}
}
