package app

import (
	"context"
	"fmt"
	"strconv"
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
	resp := &dto.APIKeyValidateResponse{}

	apiKey, err := u.apiKeyRepo.FindByKey(ctx, req.Key)
	if err != nil {
		if err == errors.ErrNotFound {
			err = errors.ErrAPIKeyNotFound
		}
		return resp, err
	}

	if !apiKey.IsActive() {
		return resp, errors.ErrAPIKeyNotActive
	}

	if apiKey.IsExpired() {
		return resp, errors.ErrAPIKeyExpired
	}

	service, err := u.serviceRepo.FindByNameAndVersion(
		ctx, req.Service, req.ServiceVersion,
	)
	if err != nil {
		if err == errors.ErrNotFound {
			err = errors.ErrServiceNotFound
		}
		return resp, err
	}

	if service.Status == enums.ServiceDeprecated {
		return resp, errors.ErrServiceDeprecated
	}
	if service.Status == enums.ServiceDeactivated {
		return resp, errors.ErrServiceDeactivated
	}

	availableRequest, err := u.environmentRepo.DecrementAvailableRequest(
		ctx, apiKey.EnvironmentID, service.ID,
	)
	if err != nil {
		return resp, errors.ErrNoAvailableRequests
	}

	requestLog := entities.RequestLog{
		APIKey:          apiKey.Key,
		ServiceID:       service.ID,
		RequestTime:     req.RequestTime,
		EnvironmentID:   apiKey.EnvironmentID,
		ExecutionStatus: enums.RequestLogPending,
	}
	if err := u.requestLog.Save(ctx, &requestLog); err != nil {
		return resp, err
	}

	var availableRequestResp string
	if availableRequest.MaxRequest == -1 {
		availableRequestResp = "unlimited"
	} else {
		availableRequestResp = strconv.Itoa(availableRequest.AvailableRequest)
	}

	println(availableRequestResp)

	return &dto.APIKeyValidateResponse{
		// Valid:     true,
		// RequestID: requestLog.ID,
		// AvailableRequest: availableRequestResp,
	}, nil
}

func (u *APIKeyUseCase) ValidateAndReserve(
	ctx context.Context, req *dto.APIKeyValidate,
) (*dto.APIKeyValidateReserveResponse, *errors.Error) {
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
) (*dto.APIKeyValidateReserveResponse, *entities.RequestLog, *entities.Reservation, *errors.Error) {
	// Validate Key, must be present in api_key entity
	apiKey, err := u.apiKeyRepo.FindByKey(ctx, req.Key)
	if err != nil {
		if err == errors.ErrNotFound {
			message := "API Key not found"
			return &dto.APIKeyValidateReserveResponse{
					Valid:   false,
					Message: message,
					Code:    enums.ReserveExecutionStatusKeyNotFound,
				},
				&entities.RequestLog{
					APIKey:          req.Key,
					RequestTime:     req.RequestTime,
					ExecutionStatus: enums.RequestLogUnauthorized,
					Message:         message,
				}, nil, nil
		}
		return nil, nil, nil, err
	}

	// Key must be active
	if !apiKey.IsActive() {
		message := "API Key is not active"
		return &dto.APIKeyValidateReserveResponse{
				Valid:   false,
				Message: message,
				Code:    enums.ReserveExecutionStatusDeactivatedKey,
			}, &entities.RequestLog{
				APIKey:          apiKey.Key,
				RequestTime:     req.RequestTime,
				EnvironmentID:   apiKey.EnvironmentID,
				ExecutionStatus: enums.RequestLogUnauthorized,
				Message:         message,
			}, nil, nil
	}

	// Expired keys are not accepted
	if apiKey.IsExpired() {
		message := "API Key expired"
		return &dto.APIKeyValidateReserveResponse{
				Valid:   false,
				Message: message,
				Code:    enums.ReserveExecutionStatusExpiredKey,
			},
			&entities.RequestLog{
				APIKey:          apiKey.Key,
				RequestTime:     req.RequestTime,
				EnvironmentID:   apiKey.EnvironmentID,
				ExecutionStatus: enums.RequestLogUnauthorized,
				Message:         message,
			}, nil, nil
	}

	environmentIsActive, err := u.environmentRepo.IsActive(
		ctx, apiKey.EnvironmentID)
	if err != nil {
		return nil, nil, nil, err
	}
	if !environmentIsActive {
		message := "Environment is deactivated"
		return &dto.APIKeyValidateReserveResponse{
				Valid:   false,
				Message: message,
				Code:    enums.ReserveExecutionStatusDeactivatedEnvironment,
			},
			&entities.RequestLog{
				APIKey:          apiKey.Key,
				RequestTime:     req.RequestTime,
				EnvironmentID:   apiKey.EnvironmentID,
				ExecutionStatus: enums.RequestLogFailed,
				Message:         message,
			}, nil, nil
	}

	// Service in that version must be exist in the service entity
	service, err := u.serviceRepo.FindByNameAndVersion(
		ctx, req.Service, req.ServiceVersion,
	)
	if err != nil {
		if err == errors.ErrNotFound {
			message := "Service not found"
			return &dto.APIKeyValidateReserveResponse{
					Valid:   false,
					Message: message,
					Code:    enums.ReserveExecutionStatusServiceNotFound,
				},
				&entities.RequestLog{
					APIKey:          apiKey.Key,
					RequestTime:     req.RequestTime,
					EnvironmentID:   apiKey.EnvironmentID,
					ExecutionStatus: enums.RequestLogFailed,
					Message:         message,
				}, nil, nil
		}
		return nil, nil, nil, err
	}

	// Service must be active
	if service.Status == enums.ServiceDeprecated {
		message := "Service is deprecated"
		return &dto.APIKeyValidateReserveResponse{
				Valid:   false,
				Message: message,
				Code:    enums.ReserveExecutionStatusDeprecatedService,
			},
			&entities.RequestLog{
				APIKey:          apiKey.Key,
				ServiceID:       service.ID,
				RequestTime:     req.RequestTime,
				EnvironmentID:   apiKey.EnvironmentID,
				ExecutionStatus: enums.RequestLogFailed,
				Message:         message,
			}, nil, nil
	}
	if service.Status == enums.ServiceDeactivated {
		message := "Service is deactivated"
		return &dto.APIKeyValidateReserveResponse{
				Valid:   false,
				Message: message,
				Code:    enums.ReserveExecutionStatusDeactivatedService,
			},
			&entities.RequestLog{
				APIKey:          apiKey.Key,
				ServiceID:       service.ID,
				RequestTime:     req.RequestTime,
				EnvironmentID:   apiKey.EnvironmentID,
				ExecutionStatus: enums.RequestLogFailed,
				Message:         message,
			}, nil, nil
	}

	// Decrement available_request or identify unlimited request
	availableRequest, err := u.environmentRepo.DecrementAvailableRequest(
		ctx, apiKey.EnvironmentID, service.ID,
	)
	if err != nil {
		/* When available_request isn't possible to decrement it must check
		the active reservations for this service in the environment no matter
		what key you use.
		*/
		currentReservations, err := u.reservationRepo.CountByEnvironmentAndService(
			ctx, apiKey.EnvironmentID, service.ID,
		)
		if err != nil {
			return nil, nil, nil, err
		}

		if currentReservations == 0 {
			message := "No available requests"
			return &dto.APIKeyValidateReserveResponse{
					Valid:   false,
					Message: message,
					Code:    enums.ReserveExecutionStatusExceededRequests,
				},
				&entities.RequestLog{
					APIKey:          apiKey.Key,
					ServiceID:       service.ID,
					RequestTime:     req.RequestTime,
					EnvironmentID:   apiKey.EnvironmentID,
					ExecutionStatus: enums.RequestLogFailed,
					Message:         message,
				}, nil, nil
		}

		if currentReservations > 0 {
			message := fmt.Sprintf("(%d) reservations are being processed and no request is available, please try again later", currentReservations)
			return &dto.APIKeyValidateReserveResponse{
					Valid:   false,
					Message: message,
					Code:    enums.ReserveExecutionStatusActiveReservations,
				},
				&entities.RequestLog{
					APIKey:          apiKey.Key,
					ServiceID:       service.ID,
					RequestTime:     req.RequestTime,
					EnvironmentID:   apiKey.EnvironmentID,
					ExecutionStatus: enums.RequestLogFailed,
					Message:         message,
				}, nil, nil
		}
	}

	// With an successful key use, then update last_used in api_key entity
	if err := u.apiKeyRepo.UpdateLastUsed(ctx, apiKey.Key); err != nil {
		return nil, nil, nil, err
	}
	return &dto.APIKeyValidateReserveResponse{
			AvailableRequest: availableRequest.AvailableRequest,
			Valid:            true,
		}, &entities.RequestLog{
			APIKey:          apiKey.Key,
			ServiceID:       service.ID,
			RequestTime:     req.RequestTime,
			EnvironmentID:   apiKey.EnvironmentID,
			ExecutionStatus: enums.RequestLogPending,
		}, &entities.Reservation{
			APIKey:        apiKey.Key,
			ServiceID:     service.ID,
			EnvironmentID: apiKey.EnvironmentID,
			RequestTime:   req.RequestTime,
			ExpiresAt:     time.Now().Add(12 * time.Hour),
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
