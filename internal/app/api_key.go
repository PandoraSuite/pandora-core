package app

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

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
) (*dto.RequestLogResponse, error) {
	service, err := u.serviceRepo.FindByNameAndVersion(
		ctx, req.ServiceName, req.ServiceVersion,
	)
	if err != nil {
		return nil, err
	}

	if service.Status != enums.ServiceActive {
		return nil, fmt.Errorf("service is %s", service.Status)
	}

	apiKey, err := u.apiKeyRepo.FindByKey(ctx, req.Key)
	if err != nil {
		return nil, err

	}

	err = u.environmentService.DecrementAvailableRequest(
		ctx, apiKey.EnvironmentID, service.ID,
	)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return &dto.RequestLogResponse{
		ID:              requestLog.ID,
		APIKey:          requestLog.APIKey,
		ServiceID:       requestLog.ServiceID,
		RequestTime:     requestLog.RequestTime,
		EnvironmentID:   requestLog.EnvironmentID,
		ExecutionStatus: requestLog.ExecutionStatus,
		CreatedAt:       requestLog.CreatedAt,
	}, nil
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
