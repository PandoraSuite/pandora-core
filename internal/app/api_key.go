package app

import (
	"context"
	"errors"
	"strconv"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	domainErr "github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/ports/outbound"
)

type APIKeyUseCase struct {
	apiKeyRepo         outbound.APIKeyPort
	requestLog         outbound.RequestLogPort
	serviceRepo        outbound.ServiceFindPort
	environmentService outbound.EnvironmentServiceQuotaPort
}

func (u *APIKeyUseCase) ValidateAndConsume(
	ctx context.Context, req *dto.APIKeyValidateAndConsume,
) (*dto.APIKeyValidateResponse, error) {
	resp := &dto.APIKeyValidateResponse{Valid: false}

	apiKey, err := u.apiKeyRepo.FindByKey(ctx, req.Key)
	if err != nil {
		if errors.Is(err, domainErr.ErrNotFound) {
			err = domainErr.ErrAPIKeyNotFound
		}
		return resp, err
	}

	if apiKey.IsExpired() {
		return resp, domainErr.ErrAPIKeyExpired
	}

	service, err := u.serviceRepo.FindByNameAndVersion(
		ctx, req.ServiceName, req.ServiceVersion,
	)
	if err != nil {
		if errors.Is(err, domainErr.ErrNotFound) {
			err = domainErr.ErrServiceNotFound
		}
		return resp, err
	}

	if service.Status == enums.ServiceDeprecated {
		return resp, domainErr.ErrServiceDeprecated
	}
	if service.Status == enums.ServiceDeactivated {
		return resp, domainErr.ErrServiceDeactivated
	}

	environmentService, err := u.environmentService.DecrementAvailableRequest(
		ctx, apiKey.EnvironmentID, service.ID,
	)
	if err != nil {
		return resp, domainErr.ErrNoAvailableRequests
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
		return resp, err
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
	}, nil
}

func (u *APIKeyUseCase) GetAPIKeysByEnvironment(
	ctx context.Context, environmentID int,
) ([]*dto.APIKeyResponse, error) {
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
) (*dto.APIKeyResponse, error) {
	apiKey := &entities.APIKey{
		Status:        enums.APIKeyActive,
		ExpiresAt:     req.ExpiresAt,
		EnvironmentID: req.EnvironmentID,
	}

	for {
		err := apiKey.GenerateKey()
		if err != nil {
			return nil, domainErr.ErrAPIKeyGenerationFailed
		}

		exists, err := u.apiKeyRepo.Exists(ctx, apiKey.Key)
		if err != nil {
			return nil, err
		}

		if !exists {
			break
		}
	}

	apiKey, err := u.apiKeyRepo.Save(ctx, apiKey)
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

func NewAPIKeyUseCase(
	apiKeyRepo outbound.APIKeyPort,
	requestLog outbound.RequestLogPort,
	serviceRepo outbound.ServiceFindPort,
	environmentService outbound.EnvironmentServiceQuotaPort,
) *APIKeyUseCase {
	return &APIKeyUseCase{
		apiKeyRepo:         apiKeyRepo,
		requestLog:         requestLog,
		serviceRepo:        serviceRepo,
		environmentService: environmentService,
	}
}
