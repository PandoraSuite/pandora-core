package validateconsume

import (
	"context"
	"fmt"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, req *dto.APIKeyValidate) (*dto.APIKeyValidateResponse, errors.Error)
}

type useCase struct {
	validator validator.Validator

	apiKeyRepo      APIKeyRepository
	serviceRepo     ServiceRepository
	requestRepo     RequestRepository
	reservationRepo ReservationRepository
	environmentRepo EnvironmentRepository
}

func (uc *useCase) Execute(
	ctx context.Context, req *dto.APIKeyValidate,
) (*dto.APIKeyValidateResponse, errors.Error) {
	validate, request, err := uc.validateAndConsume(ctx, req)
	if err != nil {
		return nil, err
	}
	// Create a request_log as a initial point, then start_point is the register id
	if err := uc.requestRepo.CreateAsInitialPoint(ctx, request); err != nil {
		return nil, err
	}
	// Return request created
	validate.RequestID = request.ID

	return validate, nil
}

func (uc *useCase) validateAndConsume(
	ctx context.Context, req *dto.APIKeyValidate,
) (*dto.APIKeyValidateResponse, *entities.Request, errors.Error) {

	// API Key must be valid, active and not expirated.
	apiKey,
		validate, request_log, err := uc.apiKeyEnable(ctx, req)
	if apiKey == nil {
		return validate, request_log, err
	}

	// Environment must be valid, active and matched with API Key
	environment,
		validate, request_log, err := uc.environmentEnable(ctx, req, apiKey)
	if environment == nil {
		return validate, request_log, err
	}

	// Service must be valid and active
	service,
		validate, request_log, err := uc.serviceEnable(ctx, req, apiKey)
	if service == nil {
		return validate, request_log, err
	}

	// Decrement available_request or identify unlimited request
	availableRequest, err := uc.environmentRepo.DecrementAvailableRequest(
		ctx, environment.ID, service.ID,
	)
	if err != nil {
		if err.Code() == errors.CodeNotFound {
			validate, request_log, err := uc.handlerErrorNotQuotes(
				ctx, req, apiKey, environment.ID, service.ID,
			)
			return validate, request_log, err
		}
		return nil, nil, err
	}
	// With an successful key use, then update last_used in api_key entity
	if err := uc.apiKeyRepo.UpdateLastUsed(ctx, apiKey.Key); err != nil {
		return nil, nil, err
	}
	return &dto.APIKeyValidateResponse{
			AvailableRequest: availableRequest.AvailableRequest,
			Valid:            true,
		}, &entities.Request{
			APIKey:          apiKey.Key,
			ServiceID:       service.ID,
			RequestTime:     req.RequestTime,
			EnvironmentID:   environment.ID,
			ExecutionStatus: enums.RequestExecutionStatusForwarded,
		}, nil
}

func (uc *useCase) apiKeyEnable(
	ctx context.Context, req *dto.APIKeyValidate,
) (*entities.APIKey, *dto.APIKeyValidateResponse, *entities.Request, errors.Error) {
	apiKey, err := uc.apiKeyRepo.GetByKey(ctx, req.APIKey)
	if err != nil {
		if err.Code() == errors.CodeNotFound {
			message := "API Key not found"
			return nil, &dto.APIKeyValidateResponse{
					Valid:   false,
					Message: message,
					Code:    enums.ValidateStatusKeyNotFound,
				},
				&entities.Request{
					APIKey:          req.APIKey,
					RequestTime:     req.RequestTime,
					ExecutionStatus: enums.RequestExecutionStatusUnauthorized,
				}, nil
		}
		return nil, nil, nil, err
	}

	// Key must be active
	if !apiKey.IsEnabled() {
		message := "API Key is not enabled"
		return nil, &dto.APIKeyValidateResponse{
				Valid:   false,
				Message: message,
				Code:    enums.ValidateStatusDeactivatedKey,
			}, &entities.Request{
				APIKey:          apiKey.Key,
				RequestTime:     req.RequestTime,
				EnvironmentID:   apiKey.EnvironmentID,
				ExecutionStatus: enums.RequestExecutionStatusUnauthorized,
			}, nil
	}

	// Expired keys are not accepted
	if apiKey.IsExpired() {
		message := "API Key expired"
		return nil, &dto.APIKeyValidateResponse{
				Valid:   false,
				Message: message,
				Code:    enums.ValidateStatusExpiredKey,
			},
			&entities.Request{
				APIKey:          apiKey.Key,
				RequestTime:     req.RequestTime,
				EnvironmentID:   apiKey.EnvironmentID,
				ExecutionStatus: enums.RequestExecutionStatusUnauthorized,
			}, nil
	}
	return apiKey, nil, nil, nil
}

func (uc *useCase) environmentEnable(
	ctx context.Context, req *dto.APIKeyValidate, apiKey *entities.APIKey,
) (*entities.Environment, *dto.APIKeyValidateResponse, *entities.Request, errors.Error) {
	environment, err := uc.environmentRepo.GetByID(
		ctx, apiKey.EnvironmentID)
	if err != nil {
		return nil, nil, nil, err
	}
	if environment.Name != req.Environment {
		message := "API Key doesn't belong to the environment"
		return nil, &dto.APIKeyValidateResponse{
				Valid:   false,
				Message: message,
				Code:    enums.ValidateStatusInvalidEnvironmentKey,
			}, &entities.Request{
				APIKey:          apiKey.Key,
				RequestTime:     req.RequestTime,
				EnvironmentID:   environment.ID,
				ExecutionStatus: enums.RequestExecutionStatusUnauthorized,
			}, nil
	}
	if !environment.IsEnabled() {
		message := "Environment is not enabled"
		return nil, &dto.APIKeyValidateResponse{
				Valid:   false,
				Message: message,
				Code:    enums.ValidateStatusDeactivatedEnvironment,
			}, &entities.Request{
				APIKey:          apiKey.Key,
				RequestTime:     req.RequestTime,
				EnvironmentID:   environment.ID,
				ExecutionStatus: enums.RequestExecutionStatusUnauthorized,
			}, nil
	}
	return environment, nil, nil, nil
}
func (uc *useCase) serviceEnable(
	ctx context.Context, req *dto.APIKeyValidate, apiKey *entities.APIKey,
) (*entities.Service, *dto.APIKeyValidateResponse, *entities.Request, errors.Error) {
	// Service in that version must be exist in the service entity
	service, err := uc.serviceRepo.GetByNameAndVersion(
		ctx, req.Service, req.ServiceVersion,
	)
	if err != nil {
		if err.Code() == errors.CodeNotFound {
			message := "Service not found"
			return nil, &dto.APIKeyValidateResponse{
					Valid:   false,
					Message: message,
					Code:    enums.ValidateStatusServiceNotFound,
				},
				&entities.Request{
					APIKey:          apiKey.Key,
					RequestTime:     req.RequestTime,
					EnvironmentID:   apiKey.EnvironmentID,
					ExecutionStatus: enums.RequestExecutionStatusUnauthorized,
				}, nil
		}
		return nil, nil, nil, err
	}

	// Service must be active
	if !service.IsEnabled() {
		message := "Service is not enabled"
		return nil, &dto.APIKeyValidateResponse{
				Valid:   false,
				Message: message,
				Code:    enums.ReservationExecutionStatusServiceNotActive,
			},
			&entities.Request{
				APIKey:          apiKey.Key,
				ServiceID:       service.ID,
				RequestTime:     req.RequestTime,
				EnvironmentID:   apiKey.EnvironmentID,
				ExecutionStatus: enums.RequestExecutionStatusUnauthorized,
			}, nil
	}
	return service, nil, nil, nil
}

func (uc *useCase) handlerErrorNotQuotes(
	ctx context.Context,
	req *dto.APIKeyValidate,
	apiKey *entities.APIKey,
	environment_id int,
	service_id int,
) (*dto.APIKeyValidateResponse, *entities.Request, errors.Error) {
	environment_service_found, has_available_requests,
		err := uc.environmentRepo.MissingResourceDiagnosis(
		ctx, environment_id, service_id)
	if err != nil {
		return nil, nil, err
	}
	if !environment_service_found {
		message := "Service does not belong to the environment"
		return &dto.APIKeyValidateResponse{
				Valid:   false,
				Message: message,
				Code:    enums.ValidateStatusEnvironmentServiceInvalid,
			},
			&entities.Request{
				APIKey:          apiKey.Key,
				ServiceID:       service_id,
				RequestTime:     req.RequestTime,
				ExecutionStatus: enums.RequestExecutionStatusUnauthorized,
			}, nil
	}
	if !has_available_requests {
		/* When available_request isn't possible to decrement it must check
		the active reservations for this service in the environment no matter
		what key you use.
		*/
		currentReservations, err := uc.reservationRepo.CountByEnvironmentAndService(
			ctx, environment_id, service_id,
		)
		if err != nil {
			return nil, nil, err
		}

		if currentReservations == 0 {
			message := "No available requests"
			return &dto.APIKeyValidateResponse{
					Valid:   false,
					Message: message,
					Code:    enums.ValidateStatusExceededRequests,
				},
				&entities.Request{
					APIKey:          apiKey.Key,
					ServiceID:       service_id,
					RequestTime:     req.RequestTime,
					EnvironmentID:   apiKey.EnvironmentID,
					ExecutionStatus: enums.RequestExecutionStatusQuotaExceeded,
				}, nil
		}

		if currentReservations > 0 {
			message := fmt.Sprintf("(%d) reservations are being processed and no request is available, please try again later", currentReservations)
			return &dto.APIKeyValidateResponse{
					Valid:   false,
					Message: message,
					Code:    enums.ValidateStatusActiveReservations,
				},
				&entities.Request{
					APIKey:          apiKey.Key,
					ServiceID:       service_id,
					RequestTime:     req.RequestTime,
					EnvironmentID:   apiKey.EnvironmentID,
					ExecutionStatus: enums.RequestExecutionStatusQuotaExceeded,
				}, nil
		}
	}
	return nil, nil, nil
}

func NewUseCase(
	validator validator.Validator,
	apiKeyRepo APIKeyRepository,
	serviceRepo ServiceRepository,
	requestRepo RequestRepository,
	reservationRepo ReservationRepository,
	environmentRepo EnvironmentRepository,
) UseCase {
	return &useCase{
		validator:       validator,
		apiKeyRepo:      apiKeyRepo,
		serviceRepo:     serviceRepo,
		requestRepo:     requestRepo,
		reservationRepo: reservationRepo,
		environmentRepo: environmentRepo,
	}
}
