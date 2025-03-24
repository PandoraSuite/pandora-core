package app

import (
	"context"
	"strconv"

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
}

func (u *APIKeyUseCase) ValidateAndConsume(
	ctx context.Context, req *dto.APIKeyValidateAndConsume,
) (*dto.APIKeyValidateResponse, *errors.Error) {
	resp := &dto.APIKeyValidateResponse{Valid: false}

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
		ctx, req.ServiceName, req.ServiceVersion,
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
	if availableRequest.MaxRequest == 0 {
		availableRequestResp = "unlimited"
	} else {
		availableRequestResp = strconv.Itoa(availableRequest.AvailableRequest)
	}

	return &dto.APIKeyValidateResponse{
		Valid:            true,
		RequestLogID:     requestLog.ID,
		AvailableRequest: availableRequestResp,
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
) *APIKeyUseCase {
	return &APIKeyUseCase{
		apiKeyRepo:      apiKeyRepo,
		requestLog:      requestLog,
		serviceRepo:     serviceRepo,
		environmentRepo: environmentRepo,
	}
}
