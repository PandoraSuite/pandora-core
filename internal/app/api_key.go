package app

import (
	"context"
	"fmt"
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/ports/outbound"
)

type APIKeyUseCase struct {
	apiKeyRepo      outbound.APIKeyPort
	requestLog      outbound.RequestLogPort
	serviceRepo     outbound.ServiceFindPort
	environmentRepo outbound.EnvironmentPort
	reservationRepo outbound.ReservationPort
}

func (u *APIKeyUseCase) ValidateAndConsume(
	ctx context.Context, req *dto.APIKeyValidate,
) (*dto.APIKeyValidateResponse, *errors.Error) {
	validate, requestLog, err := u.validateAndConsume(ctx, req)
	if err != nil {
		return nil, err
	}
	// Create a request_log as a initial point, then start_point is the register id
	if err := u.requestLog.SaveAsInitialPoint(ctx, requestLog); err != nil {
		return nil, err
	}
	// Return request created
	validate.RequestID = requestLog.ID

	return validate, nil
}
func (u *APIKeyUseCase) validateAndConsume(
	ctx context.Context, req *dto.APIKeyValidate,
) (*dto.APIKeyValidateResponse, *entities.RequestLog, *errors.Error) {

	// API Key must be valid, active and not expirated.
	apiKey,
		validate, request_log, err := u.apiKeyEnable(ctx, req)
	if apiKey == nil {
		return validate, request_log, err
	}

	// Environment must be valid, active and matched with API Key
	environment,
		validate, request_log, err := u.environmentEnable(ctx, req, apiKey)
	if environment == nil {
		return validate, request_log, err
	}

	// Service must be valid and active
	service,
		validate, request_log, err := u.serviceEnable(ctx, req, apiKey)
	if service == nil {
		return validate, request_log, err
	}

	// Decrement available_request or identify unlimited request
	availableRequest, err := u.environmentRepo.DecrementAvailableRequest(
		ctx, environment.ID, service.ID,
	)
	if err != nil {
		if err == errors.ErrNotFound {
			validate, request_log, err := u.handlerErrorNotQuotes(
				ctx, req, apiKey, environment.ID, service.ID,
			)
			return validate, request_log, err
		}
		return nil, nil, err
	}
	// With an successful key use, then update last_used in api_key entity
	if err := u.apiKeyRepo.UpdateLastUsed(ctx, apiKey.Key); err != nil {
		return nil, nil, err
	}
	return &dto.APIKeyValidateResponse{
			AvailableRequest: availableRequest.AvailableRequest,
			Valid:            true,
		}, &entities.RequestLog{
			APIKey:          apiKey.Key,
			ServiceID:       service.ID,
			RequestTime:     req.RequestTime,
			EnvironmentID:   environment.ID,
			ExecutionStatus: enums.RequestLogPending,
		}, nil
}

func (u *APIKeyUseCase) ValidateAndReserve(
	ctx context.Context, req *dto.APIKeyValidate,
) (*dto.APIKeyValidateResponse, *errors.Error) {
	validate, requestLog, reservation, err := u.validateAndReserve(ctx, req)
	if err != nil {
		return nil, err
	}
	// Create a request_log as a initial point, then start_point is the register id
	if err := u.requestLog.SaveAsInitialPoint(ctx, requestLog); err != nil {
		return nil, err
	}
	// Return request created
	validate.RequestID = requestLog.ID

	// When process was valid then create an active reservation up to twelve hours later
	if reservation != nil {
		reservation.StartRequestID = requestLog.ID
		if err := u.reservationRepo.Save(ctx, reservation); err != nil {
			return nil, err
		}
		validate.ReservationID = reservation.ID
	}

	return validate, nil
}

func (u *APIKeyUseCase) validateAndReserve(
	ctx context.Context, req *dto.APIKeyValidate,
) (*dto.APIKeyValidateResponse, *entities.RequestLog, *entities.Reservation, *errors.Error) {
	// API Key must be valid, active and not expirated.
	apiKey,
		validate, request_log, err := u.apiKeyEnable(ctx, req)
	if apiKey == nil {
		return validate, request_log, nil, err
	}

	// Environment must be valid, active and matched with API Key
	environment,
		validate, request_log, err := u.environmentEnable(ctx, req, apiKey)
	if environment == nil {
		return validate, request_log, nil, err
	}

	// Service must be valid and active
	service,
		validate, request_log, err := u.serviceEnable(ctx, req, apiKey)
	if service == nil {
		return validate, request_log, nil, err
	}

	// Decrement available_request or identify unlimited request
	availableRequest, err := u.environmentRepo.DecrementAvailableRequest(
		ctx, environment.ID, service.ID,
	)
	if err != nil {
		if err == errors.ErrNotFound {
			validate, request_log, err := u.handlerErrorNotQuotes(
				ctx, req, apiKey, environment.ID, service.ID,
			)
			return validate, request_log, nil, err
		}
		return nil, nil, nil, err
	}

	// With an successful key use, then update last_used in api_key entity
	if err := u.apiKeyRepo.UpdateLastUsed(ctx, apiKey.Key); err != nil {
		return nil, nil, nil, err
	}
	return &dto.APIKeyValidateResponse{
			AvailableRequest: availableRequest.AvailableRequest,
			Valid:            true,
		}, &entities.RequestLog{
			APIKey:          apiKey.Key,
			ServiceID:       service.ID,
			RequestTime:     req.RequestTime,
			EnvironmentID:   environment.ID,
			ExecutionStatus: enums.RequestLogPending,
		}, &entities.Reservation{
			APIKey:        apiKey.Key,
			ServiceID:     service.ID,
			EnvironmentID: environment.ID,
			RequestTime:   req.RequestTime,
			ExpiresAt:     time.Now().Add(12 * time.Hour),
		}, nil
}

func (u *APIKeyUseCase) apiKeyEnable(
	ctx context.Context, req *dto.APIKeyValidate,
) (*entities.APIKey, *dto.APIKeyValidateResponse, *entities.RequestLog, *errors.Error) {
	apiKey, err := u.apiKeyRepo.FindByKey(ctx, req.Key)
	if err != nil {
		if err == errors.ErrNotFound {
			message := "API Key not found"
			return nil, &dto.APIKeyValidateResponse{
					Valid:   false,
					Message: message,
					Code:    enums.ValidateStatusKeyNotFound,
				},
				&entities.RequestLog{
					APIKey:          req.Key,
					RequestTime:     req.RequestTime,
					ExecutionStatus: enums.RequestLogUnauthorized,
					Message:         message,
				}, nil
		}
		return nil, nil, nil, err
	}

	// Key must be active
	if !apiKey.IsActive() {
		message := "API Key is not active"
		return nil, &dto.APIKeyValidateResponse{
				Valid:   false,
				Message: message,
				Code:    enums.ValidateStatusDeactivatedKey,
			}, &entities.RequestLog{
				APIKey:          apiKey.Key,
				RequestTime:     req.RequestTime,
				EnvironmentID:   apiKey.EnvironmentID,
				ExecutionStatus: enums.RequestLogUnauthorized,
				Message:         message,
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
			&entities.RequestLog{
				APIKey:          apiKey.Key,
				RequestTime:     req.RequestTime,
				EnvironmentID:   apiKey.EnvironmentID,
				ExecutionStatus: enums.RequestLogUnauthorized,
				Message:         message,
			}, nil
	}
	return apiKey, nil, nil, nil
}

func (u *APIKeyUseCase) environmentEnable(
	ctx context.Context, req *dto.APIKeyValidate, apiKey *entities.APIKey,
) (*entities.Environment, *dto.APIKeyValidateResponse, *entities.RequestLog, *errors.Error) {
	environment, err := u.environmentRepo.FindByID(
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
			}, &entities.RequestLog{
				APIKey:          apiKey.Key,
				RequestTime:     req.RequestTime,
				EnvironmentID:   environment.ID,
				ExecutionStatus: enums.RequestLogUnauthorized,
				Message:         message,
			}, nil
	}
	if !environment.IsActive() {
		message := "Environment is not active"
		return nil, &dto.APIKeyValidateResponse{
				Valid:   false,
				Message: message,
				Code:    enums.ValidateStatusDeactivatedEnvironment,
			}, &entities.RequestLog{
				APIKey:          apiKey.Key,
				RequestTime:     req.RequestTime,
				EnvironmentID:   environment.ID,
				ExecutionStatus: enums.RequestLogUnauthorized,
				Message:         message,
			}, nil
	}
	return environment, nil, nil, nil
}
func (u *APIKeyUseCase) serviceEnable(
	ctx context.Context, req *dto.APIKeyValidate, apiKey *entities.APIKey,
) (*entities.Service, *dto.APIKeyValidateResponse, *entities.RequestLog, *errors.Error) {
	// Service in that version must be exist in the service entity
	service, err := u.serviceRepo.FindByNameAndVersion(
		ctx, req.Service, req.ServiceVersion,
	)
	if err != nil {
		if err == errors.ErrNotFound {
			message := "Service not found"
			return nil, &dto.APIKeyValidateResponse{
					Valid:   false,
					Message: message,
					Code:    enums.ValidateStatusServiceNotFound,
				},
				&entities.RequestLog{
					APIKey:          apiKey.Key,
					RequestTime:     req.RequestTime,
					EnvironmentID:   apiKey.EnvironmentID,
					ExecutionStatus: enums.RequestLogUnauthorized,
					Message:         message,
				}, nil
		}
		return nil, nil, nil, err
	}

	// Service must be active
	if !service.IsActive() {
		message := "Service is not active"
		return nil, &dto.APIKeyValidateResponse{
				Valid:   false,
				Message: message,
				Code:    enums.ReservationExecutionStatusServiceNotActive,
			},
			&entities.RequestLog{
				APIKey:          apiKey.Key,
				ServiceID:       service.ID,
				RequestTime:     req.RequestTime,
				EnvironmentID:   apiKey.EnvironmentID,
				ExecutionStatus: enums.RequestLogUnauthorized,
				Message:         message,
			}, nil
	}
	return service, nil, nil, nil
}

func (u *APIKeyUseCase) handlerErrorNotQuotes(
	ctx context.Context,
	req *dto.APIKeyValidate,
	apiKey *entities.APIKey,
	environment_id int,
	service_id int,
) (*dto.APIKeyValidateResponse, *entities.RequestLog, *errors.Error) {
	environment_service_found, has_available_requests,
		err := u.environmentRepo.MissingResourceDiagnosis(
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
			&entities.RequestLog{
				APIKey:          apiKey.Key,
				ServiceID:       service_id,
				RequestTime:     req.RequestTime,
				ExecutionStatus: enums.RequestLogUnauthorized,
				Message:         message,
			}, nil
	}
	if !has_available_requests {
		/* When available_request isn't possible to decrement it must check
		the active reservations for this service in the environment no matter
		what key you use.
		*/
		currentReservations, err := u.reservationRepo.CountByEnvironmentAndService(
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
				&entities.RequestLog{
					APIKey:          apiKey.Key,
					ServiceID:       service_id,
					RequestTime:     req.RequestTime,
					EnvironmentID:   apiKey.EnvironmentID,
					ExecutionStatus: enums.RequestLogFailed,
					Message:         message,
				}, nil
		}

		if currentReservations > 0 {
			message := fmt.Sprintf("(%d) reservations are being processed and no request is available, please try again later", currentReservations)
			return &dto.APIKeyValidateResponse{
					Valid:   false,
					Message: message,
					Code:    enums.ValidateStatusActiveReservations,
				},
				&entities.RequestLog{
					APIKey:          apiKey.Key,
					ServiceID:       service_id,
					RequestTime:     req.RequestTime,
					EnvironmentID:   apiKey.EnvironmentID,
					ExecutionStatus: enums.RequestLogFailed,
					Message:         message,
				}, nil
		}
	}
	return nil, nil, nil
}

func (u *APIKeyUseCase) ValidateWithReservation(
	ctx context.Context, req *dto.APIKeyValidateReserve,
) (*dto.APIKeyValidateReservationResponse, *errors.Error) {
	validate, requestLog, err := u.validateWithReservation(ctx, req)
	if err != nil {
		return nil, err
	}
	// Create a request_log, then start_point is the start_request_id of reservation active
	if err := u.requestLog.Save(ctx, requestLog); err != nil {
		return nil, err
	}
	// Return request created
	validate.RequestID = requestLog.ID

	return validate, nil
}

func (u *APIKeyUseCase) validateWithReservation(
	ctx context.Context, req *dto.APIKeyValidateReserve,
) (*dto.APIKeyValidateReservationResponse, *entities.RequestLog, *errors.Error) {
	reservationFlow, err := u.reservationRepo.FindByIDWithDetails(ctx, req.ReservationID)
	if err != nil {
		if err == errors.ErrNotFound {
			message := "Reservation not found"
			return &dto.APIKeyValidateReservationResponse{
					Valid:   false,
					Message: message,
					Code:    enums.ReservationExecutionStatusNotFound,
				}, &entities.RequestLog{
					APIKey:          req.Key,
					RequestTime:     req.RequestTime,
					ExecutionStatus: enums.RequestLogFailed,
					Message:         message,
				}, nil
		}
		return nil, nil, err
	}
	if req.Key != reservationFlow.APIKey {
		message := "API Key invalid for the reservation"
		return &dto.APIKeyValidateReservationResponse{
				Valid:   false,
				Message: message,
				Code:    enums.ValidateStatusInvalidKey,
			},
			&entities.RequestLog{
				APIKey:          req.Key,
				ServiceID:       reservationFlow.ServiceID,
				RequestTime:     req.RequestTime,
				EnvironmentID:   reservationFlow.EnvironmentID,
				StartPoint:      reservationFlow.StartRequestID,
				ExecutionStatus: enums.RequestLogFailed,
				Message:         message,
			}, nil
	}
	if req.Service != reservationFlow.ServiceName {
		message := "Service invalid for the reservation"
		return &dto.APIKeyValidateReservationResponse{
				Valid:   false,
				Message: message,
				Code:    enums.ReservationExecutionStatusInvalidService,
			},
			&entities.RequestLog{
				APIKey:          reservationFlow.APIKey,
				ServiceID:       reservationFlow.ServiceID,
				RequestTime:     req.RequestTime,
				EnvironmentID:   reservationFlow.EnvironmentID,
				StartPoint:      reservationFlow.StartRequestID,
				ExecutionStatus: enums.RequestLogFailed,
				Message:         message,
			}, nil
	}

	if req.ServiceVersion != reservationFlow.ServiceVersion {
		message := "Service version invalid for the reservation"
		return &dto.APIKeyValidateReservationResponse{
				Valid:   false,
				Message: message,
				Code:    enums.ReservationExecutionStatusInvalidServiceVersion,
			}, &entities.RequestLog{
				APIKey:          reservationFlow.APIKey,
				ServiceID:       reservationFlow.ServiceID,
				RequestTime:     req.RequestTime,
				EnvironmentID:   reservationFlow.EnvironmentID,
				StartPoint:      reservationFlow.StartRequestID,
				ExecutionStatus: enums.RequestLogFailed,
				Message:         message,
			}, nil
	}

	if reservationFlow.ServiceStatus != enums.ServiceActive {
		message := "Service is not active"
		return &dto.APIKeyValidateReservationResponse{
				Valid:   false,
				Message: message,
				Code:    enums.ReservationExecutionStatusServiceNotActive,
			}, &entities.RequestLog{
				APIKey:          reservationFlow.APIKey,
				ServiceID:       reservationFlow.ServiceID,
				RequestTime:     req.RequestTime,
				EnvironmentID:   reservationFlow.EnvironmentID,
				StartPoint:      reservationFlow.StartRequestID,
				ExecutionStatus: enums.RequestLogFailed,
				Message:         message,
			}, nil
	}

	if req.Environment != reservationFlow.EnvironmentName {
		message := "Environment invalid for the reservation"
		return &dto.APIKeyValidateReservationResponse{
				Valid:   false,
				Message: message,
				Code:    enums.ReservationExecutionStatusInvalidEnvironment,
			}, &entities.RequestLog{
				APIKey:          reservationFlow.APIKey,
				ServiceID:       reservationFlow.ServiceID,
				RequestTime:     req.RequestTime,
				EnvironmentID:   reservationFlow.EnvironmentID,
				StartPoint:      reservationFlow.StartRequestID,
				ExecutionStatus: enums.RequestLogFailed,
				Message:         message,
			}, nil
	}

	if reservationFlow.EnvironmentStatus != enums.EnvironmentActive {
		message := "Environment is not active"
		return &dto.APIKeyValidateReservationResponse{
				Valid:   false,
				Message: message,
				Code:    enums.ReservationExecutionStatusEnvironmentNotActive,
			}, &entities.RequestLog{
				APIKey:          reservationFlow.APIKey,
				ServiceID:       reservationFlow.ServiceID,
				RequestTime:     req.RequestTime,
				EnvironmentID:   reservationFlow.EnvironmentID,
				StartPoint:      reservationFlow.StartRequestID,
				ExecutionStatus: enums.RequestLogFailed,
				Message:         message,
			}, nil
	}

	return &dto.APIKeyValidateReservationResponse{
			Valid: true,
		}, &entities.RequestLog{
			APIKey:          reservationFlow.APIKey,
			ServiceID:       reservationFlow.ServiceID,
			RequestTime:     req.RequestTime,
			EnvironmentID:   reservationFlow.EnvironmentID,
			StartPoint:      reservationFlow.StartRequestID,
			ExecutionStatus: enums.RequestLogPending,
		}, nil

}

func (u *APIKeyUseCase) GetAPIKeysByEnvironment(
	ctx context.Context, environmentID int,
) ([]*dto.APIKeyResponse, *errors.Error) {
	apiKeys, err := u.apiKeyRepo.FindByEnvironment(ctx, environmentID)
	if err != nil {
		return nil, err
	}

	apiKeysResponses := make([]*dto.APIKeyResponse, len(apiKeys))
	for i, apiKey := range apiKeys {
		apiKeysResponses[i] = &dto.APIKeyResponse{
			ID:            apiKey.ID,
			Key:           apiKey.Key,
			Status:        apiKey.Status,
			LastUsed:      apiKey.LastUsed,
			ExpiresAt:     apiKey.ExpiresAt,
			EnvironmentID: apiKey.EnvironmentID,
			CreatedAt:     apiKey.CreatedAt,
		}
	}

	return apiKeysResponses, nil
}

func (u *APIKeyUseCase) Update(
	ctx context.Context, id int, req *dto.APIKeyUpdate,
) (*dto.APIKeyResponse, *errors.Error) {
	if !req.ExpiresAt.IsZero() && req.ExpiresAt.Before(time.Now()) {
		return nil, errors.ErrAPIKeyInvalidExpiresAt
	}

	apiKey, err := u.apiKeyRepo.Update(ctx, id, req)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrAPIKeyNotFound
		}
		return nil, err
	}

	return &dto.APIKeyResponse{
		ID:            apiKey.ID,
		Key:           apiKey.Key,
		Status:        apiKey.Status,
		LastUsed:      apiKey.LastUsed,
		ExpiresAt:     apiKey.ExpiresAt,
		EnvironmentID: apiKey.EnvironmentID,
		CreatedAt:     apiKey.CreatedAt,
	}, nil
}

func (u *APIKeyUseCase) Create(
	ctx context.Context, req *dto.APIKeyCreate,
) (*dto.APIKeyResponse, *errors.Error) {
	apiKey := entities.APIKey{
		Status:        enums.APIKeyActive,
		ExpiresAt:     req.ExpiresAt,
		EnvironmentID: req.EnvironmentID,
	}

	for {
		err := apiKey.GenerateKey()
		if err != nil {
			return nil, err
		}

		exists, err := u.apiKeyRepo.Exists(ctx, apiKey.Key)
		if err != nil {
			return nil, err
		}

		if !exists {
			break
		}
	}

	if err := u.apiKeyRepo.Save(ctx, &apiKey); err != nil {
		return nil, err
	}

	return &dto.APIKeyResponse{
		ID:            apiKey.ID,
		Key:           apiKey.Key,
		Status:        apiKey.Status,
		LastUsed:      apiKey.LastUsed,
		ExpiresAt:     apiKey.ExpiresAt,
		EnvironmentID: apiKey.EnvironmentID,
		CreatedAt:     apiKey.CreatedAt,
	}, nil
}

func NewAPIKeyUseCase(
	apiKeyRepo outbound.APIKeyPort,
	requestLog outbound.RequestLogPort,
	serviceRepo outbound.ServiceFindPort,
	environmentRepo outbound.EnvironmentPort,
	reservationRepo outbound.ReservationPort,
) *APIKeyUseCase {
	return &APIKeyUseCase{
		apiKeyRepo:      apiKeyRepo,
		requestLog:      requestLog,
		serviceRepo:     serviceRepo,
		environmentRepo: environmentRepo,
		reservationRepo: reservationRepo,
	}
}
