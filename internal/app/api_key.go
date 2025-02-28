package app

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"strconv"
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/ports/outbound"
)

type APIKeyUseCase struct {
	apiKeyRepo         outbound.APIKeyRepositoryPort
	requestLog         outbound.RequestLogRepositoryPort
	serviceRepo        outbound.ServiceRepositoryPort
	environmentService outbound.EnvironmentServiceRepositoryPort
}

func (u *APIKeyUseCase) ValidateAndConsume(
	ctx context.Context, req *dto.APIKeyValidateAndConsume,
) *dto.APIKeyValidateResponse {
	apiKey, err := u.apiKeyRepo.FindByKey(ctx, req.Key)
	if err != nil {
		return &dto.APIKeyValidateResponse{
			Valid:   false,
			Message: err.Error(),
		}
	}

	if apiKey.ExpiresAt.Before(time.Now()) {
		return &dto.APIKeyValidateResponse{
			Valid:   false,
			Message: "api key has expired",
		}
	}

	service, err := u.serviceRepo.FindByNameAndVersion(
		ctx, req.ServiceName, req.ServiceVersion,
	)
	if err != nil {
		return &dto.APIKeyValidateResponse{
			Valid:   false,
			Message: err.Error(),
		}
	}

	if service.Status == enums.ServiceDeprecated {
		return &dto.APIKeyValidateResponse{
			Valid:   false,
			Message: "service is deprecated",
		}
	}
	if service.Status == enums.ServiceDeactivated {
		return &dto.APIKeyValidateResponse{
			Valid:   false,
			Message: "service is not available",
		}
	}

	environmentService, err := u.environmentService.DecrementAvailableRequest(
		ctx, apiKey.EnvironmentID, service.ID,
	)
	if err != nil {
		return &dto.APIKeyValidateResponse{
			Valid:   false,
			Message: err.Error(),
		}
	}

	requestLog, err := u.requestLog.Save(
		ctx,
		&entities.RequestLog{
			APIKey:          apiKey.Key,
			ServiceID:       service.ID,
			RequestTime:     req.RequestTime,
			EnvironmentID:   apiKey.EnvironmentID,
			ExecutionStatus: enums.RequestLogPending,
		},
	)
	if err != nil {
		return &dto.APIKeyValidateResponse{
			Valid:   false,
			Message: err.Error(),
		}
	}

	var availableRequest string
	if environmentService.MaxRequest == 0 {
		availableRequest = "unlimited"
	} else {
		availableRequest = strconv.Itoa(environmentService.AvailableRequest)
	}

	return &dto.APIKeyValidateResponse{
		Valid:            true,
		AvailableRequest: availableRequest,
		RequestLogID:     requestLog.ID,
	}
}

func (u *APIKeyUseCase) Create(
	ctx context.Context, req *dto.APIKeyCreate,
) (*dto.APIKeyResponse, error) {
	var key string
	for {
		key, err := u.generateKey()
		if err != nil {
			return nil, err // TODO: handle error
		}

		exists, err := u.apiKeyRepo.Exists(ctx, key)
		if err != nil {
			return nil, err // TODO: handle error
		}

		if !exists {
			break
		}
	}

	apiKey, err := u.apiKeyRepo.Save(
		ctx,
		&entities.APIKey{
			Key:           key,
			Status:        enums.APIKeyActive,
			ExpiresAt:     req.ExpiresAt,
			EnvironmentID: req.EnvironmentID,
		},
	)
	if err != nil {
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

func (u *APIKeyUseCase) generateKey() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
