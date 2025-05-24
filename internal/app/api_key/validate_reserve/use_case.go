package validatereserve

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, req *dto.APIKeyValidateReserve) (*dto.APIKeyValidateReservationResponse, errors.Error)
}

type useCase struct {
	validator validator.Validator

	requestRepo     RequestRepository
	reservationRepo ReservationRepository
}

func (uc *useCase) Execute(
	ctx context.Context, req *dto.APIKeyValidateReserve,
) (*dto.APIKeyValidateReservationResponse, errors.Error) {
	validate, request, err := uc.validateWithReservation(ctx, req)
	if err != nil {
		return nil, err
	}
	// Create a request_log, then start_point is the start_request_id of reservation active
	if err := uc.requestRepo.Create(ctx, request); err != nil {
		return nil, err
	}
	// Return request created
	validate.RequestID = request.ID

	return validate, nil
}

func (uc *useCase) validateWithReservation(
	ctx context.Context, req *dto.APIKeyValidateReserve,
) (*dto.APIKeyValidateReservationResponse, *entities.Request, errors.Error) {
	reservationFlow, err := uc.reservationRepo.GetByIDWithDetails(ctx, req.ReservationID)
	if err != nil {
		if err.Code() == errors.CodeNotFound {
			message := "Reservation not found"
			return &dto.APIKeyValidateReservationResponse{
					Valid:   false,
					Message: message,
					Code:    enums.ReservationExecutionStatusNotFound,
				}, &entities.Request{
					APIKey:          req.APIKey,
					RequestTime:     req.RequestTime,
					ExecutionStatus: enums.RequestExecutionStatusQuotaExceeded,
				}, nil
		}
		return nil, nil, err
	}
	if req.APIKey != reservationFlow.APIKey {
		message := "API Key invalid for the reservation"
		return &dto.APIKeyValidateReservationResponse{
				Valid:   false,
				Message: message,
				Code:    enums.ValidateStatusInvalidKey,
			},
			&entities.Request{
				APIKey:          req.APIKey,
				ServiceID:       reservationFlow.ServiceID,
				RequestTime:     req.RequestTime,
				EnvironmentID:   reservationFlow.EnvironmentID,
				StartPoint:      reservationFlow.StartRequestID,
				ExecutionStatus: enums.RequestExecutionStatusQuotaExceeded,
			}, nil
	}
	if req.Service != reservationFlow.ServiceName {
		message := "Service invalid for the reservation"
		return &dto.APIKeyValidateReservationResponse{
				Valid:   false,
				Message: message,
				Code:    enums.ReservationExecutionStatusInvalidService,
			},
			&entities.Request{
				APIKey:          reservationFlow.APIKey,
				ServiceID:       reservationFlow.ServiceID,
				RequestTime:     req.RequestTime,
				EnvironmentID:   reservationFlow.EnvironmentID,
				StartPoint:      reservationFlow.StartRequestID,
				ExecutionStatus: enums.RequestExecutionStatusQuotaExceeded,
			}, nil
	}

	if req.ServiceVersion != reservationFlow.ServiceVersion {
		message := "Service version invalid for the reservation"
		return &dto.APIKeyValidateReservationResponse{
				Valid:   false,
				Message: message,
				Code:    enums.ReservationExecutionStatusInvalidServiceVersion,
			}, &entities.Request{
				APIKey:          reservationFlow.APIKey,
				ServiceID:       reservationFlow.ServiceID,
				RequestTime:     req.RequestTime,
				EnvironmentID:   reservationFlow.EnvironmentID,
				StartPoint:      reservationFlow.StartRequestID,
				ExecutionStatus: enums.RequestExecutionStatusQuotaExceeded,
			}, nil
	}

	if reservationFlow.ServiceStatus != enums.ServiceActive {
		message := "Service is not active"
		return &dto.APIKeyValidateReservationResponse{
				Valid:   false,
				Message: message,
				Code:    enums.ReservationExecutionStatusServiceNotActive,
			}, &entities.Request{
				APIKey:          reservationFlow.APIKey,
				ServiceID:       reservationFlow.ServiceID,
				RequestTime:     req.RequestTime,
				EnvironmentID:   reservationFlow.EnvironmentID,
				StartPoint:      reservationFlow.StartRequestID,
				ExecutionStatus: enums.RequestExecutionStatusQuotaExceeded,
			}, nil
	}

	if req.Environment != reservationFlow.EnvironmentName {
		message := "Environment invalid for the reservation"
		return &dto.APIKeyValidateReservationResponse{
				Valid:   false,
				Message: message,
				Code:    enums.ReservationExecutionStatusInvalidEnvironment,
			}, &entities.Request{
				APIKey:          reservationFlow.APIKey,
				ServiceID:       reservationFlow.ServiceID,
				RequestTime:     req.RequestTime,
				EnvironmentID:   reservationFlow.EnvironmentID,
				StartPoint:      reservationFlow.StartRequestID,
				ExecutionStatus: enums.RequestExecutionStatusQuotaExceeded,
			}, nil
	}

	if reservationFlow.EnvironmentStatus != enums.EnvironmentActive {
		message := "Environment is not active"
		return &dto.APIKeyValidateReservationResponse{
				Valid:   false,
				Message: message,
				Code:    enums.ReservationExecutionStatusEnvironmentNotActive,
			}, &entities.Request{
				APIKey:          reservationFlow.APIKey,
				ServiceID:       reservationFlow.ServiceID,
				RequestTime:     req.RequestTime,
				EnvironmentID:   reservationFlow.EnvironmentID,
				StartPoint:      reservationFlow.StartRequestID,
				ExecutionStatus: enums.RequestExecutionStatusQuotaExceeded,
			}, nil
	}

	return &dto.APIKeyValidateReservationResponse{
			Valid: true,
		}, &entities.Request{
			APIKey:          reservationFlow.APIKey,
			ServiceID:       reservationFlow.ServiceID,
			RequestTime:     req.RequestTime,
			EnvironmentID:   reservationFlow.EnvironmentID,
			StartPoint:      reservationFlow.StartRequestID,
			ExecutionStatus: enums.RequestExecutionStatusForwarded,
		}, nil
}

func NewUseCase(
	validator validator.Validator,
	requestRepo RequestRepository,
	reservationRepo ReservationRepository,
) UseCase {
	return &useCase{
		requestRepo:     requestRepo,
		reservationRepo: reservationRepo,
		validator:       validator,
	}
}
