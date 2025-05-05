package listapikey

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, id int) ([]*dto.APIKeyResponse, errors.Error)
}

type useCase struct {
	validator validator.Validator

	apiKeyRepo      APIKeyRepository
	environmentRepo EnvironmentRepository
}

func (uc *useCase) Execute(
	ctx context.Context, id int,
) ([]*dto.APIKeyResponse, errors.Error) {
	if err := uc.validateID(id); err != nil {
		return nil, err
	}

	exists, err := uc.environmentRepo.Exists(ctx, id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.NewEntityNotFound(
			"Environment",
			"environment not found",
			map[string]any{"id": id},
		)
	}

	apiKeys, err := uc.apiKeyRepo.ListByEnvironment(ctx, id)
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

func (uc *useCase) validateID(id int) errors.Error {
	return uc.validator.ValidateVariable(
		id,
		"id",
		"required,gt=0",
		map[string]string{
			"required": "id is required",
			"gt":       "id must be greater than 0",
		},
	)
}

func NewUseCase(
	validator validator.Validator,
	apiKeyRepo APIKeyRepository,
	environmentRepo EnvironmentRepository,
) UseCase {
	return &useCase{
		validator:       validator,
		apiKeyRepo:      apiKeyRepo,
		environmentRepo: environmentRepo,
	}
}
