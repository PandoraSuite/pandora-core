package app

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/ports/outbound"
)

type APIKeyUseCase struct {
	apiKeyRepo outbound.APIKeyRepositoryPort
}

func (s *APIKeyUseCase) Create(
	ctx context.Context, req *dto.APIKeyCreate,
) (*dto.APIKeyResponse, error) {
	var key string
	for {
		key, err := s.generateKey()
		if err != nil {
			return nil, err // TODO: handle error
		}

		exists, err := s.apiKeyRepo.Exists(ctx, key)
		if err != nil {
			return nil, err // TODO: handle error
		}

		if !exists {
			break
		}
	}

	apiKey, err := s.apiKeyRepo.Save(
		ctx,
		&entities.APIKey{
			Key:           key,
			Status:        "active",
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
		Status:        dto.APIKeyStatus(apiKey.Status),
		ExpiresAt:     apiKey.ExpiresAt,
		EnvironmentID: apiKey.EnvironmentID,
		CreatedAt:     apiKey.CreatedAt,
	}, nil
}

func (s *APIKeyUseCase) generateKey() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
