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
	// Validate Key, must be present in api_key entity
	apiKey, err := u.apiKeyRepo.FindByKey(ctx, req.Key)
	if err != nil {
		if err == errors.ErrNotFound {
			return &dto.APIKeyValidateReserveResponse{
				Valid:   false,
				Message: "API Key not found",
				Code:    enums.ReserveExecutionStatusKeyNotFound,
			}, nil
		}
		return nil, err
	}

	// Key must be active
	if !apiKey.IsActive() {
		return &dto.APIKeyValidateReserveResponse{
			Valid:   false,
			Message: "API Key is not active",
			Code:    enums.ReserveExecutionStatusDeactivatedKey,
		}, nil
	}

	// Expired keys are not accepted
	if apiKey.IsExpired() {
		return &dto.APIKeyValidateReserveResponse{
			Valid:   false,
			Message: "API Key expired",
			Code:    enums.ReserveExecutionStatusExpiredKey,
		}, nil
	}

	// Service in that version must be exist in the service entity
	service, err := u.serviceRepo.FindByNameAndVersion(
		ctx, req.Service, req.ServiceVersion,
	)
	if err != nil {
		if err == errors.ErrNotFound {
			return &dto.APIKeyValidateReserveResponse{
				Valid:   false,
				Message: "Service not found",
				Code:    enums.ReserveExecutionStatusServiceNotFound,
			}, nil
		}
		return nil, err
	}

	// Service must be active
	if service.Status == enums.ServiceDeprecated {
		return &dto.APIKeyValidateReserveResponse{
			Valid:   false,
			Message: "Service is deprecated",
			Code:    enums.ReserveExecutionStatusDeprecatedService,
		}, nil
	}
	if service.Status == enums.ServiceDeactivated {
		return &dto.APIKeyValidateReserveResponse{
			Valid:   false,
			Message: "Service is deactivated",
			Code:    enums.ReserveExecutionStatusDeactivatedService,
		}, nil
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
		currentReservations, err := u.reservationRepo.CountReservationsByFields(
			ctx, apiKey.EnvironmentID, service.ID,
		)
		if err != nil {
			return nil, err
		}

		if currentReservations == 0 {
			return &dto.APIKeyValidateReserveResponse{
				Valid:   false,
				Message: "No available requests",
				Code:    enums.ReserveExecutionStatusExceededRequests,
			}, nil
		}

		if currentReservations > 0 {
			return &dto.APIKeyValidateReserveResponse{
				Valid:   false,
				Message: fmt.Sprintf("(%d) reservations are being processed and no request is available, please try again later", currentReservations),
				Code:    enums.ReserveExecutionStatusActiveReservations,
			}, nil
		}
	}

	// Create an active reservation up to twelve hours later
	expiresAtTime := time.Now().Add(12 * time.Hour)
	reservation := entities.Reservation{
		APIKey:        apiKey.Key,
		ServiceID:     service.ID,
		EnvironmentID: apiKey.EnvironmentID,
		RequestTime:   req.RequestTime,
		ExpiresAt:     expiresAtTime,
	}
	if err := u.reservationRepo.Save(ctx, &reservation); err != nil {
		return nil, err
	}
	requestLog := entities.RequestLog{
		APIKey:          apiKey.Key,
		ServiceID:       service.ID,
		RequestTime:     req.RequestTime,
		EnvironmentID:   apiKey.EnvironmentID,
		ExecutionStatus: enums.RequestLogPending,
	}
	if err := u.requestLog.SaveAsInitialPoint(ctx, &requestLog); err != nil {
		return nil, err
	}
	// With an successful key use, then update last_used in api_key entity
	if err := u.apiKeyRepo.UpdateLastUsed(ctx, apiKey.Key); err != nil {
		return nil, err
	}

	return &dto.APIKeyValidateReserveResponse{
		ReservationID:    reservation.ID,
		AvailableRequest: availableRequest.AvailableRequest,
		Valid:            true,
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
