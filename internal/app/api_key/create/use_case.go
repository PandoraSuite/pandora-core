package create

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, req *dto.APIKeyCreate) (*dto.APIKeyResponse, errors.Error)
}

type useCase struct {
	validator validator.Validator

	apiKeyRepo APIKeyRepository
}

func (uc *useCase) Execute(
	ctx context.Context, req *dto.APIKeyCreate,
) (*dto.APIKeyResponse, errors.Error) {
	if err := uc.validateReq(req); err != nil {
		return nil, err
	}

	apiKey := &entities.APIKey{
		Status:        enums.APIKeyStatusEnabled,
		ExpiresAt:     req.ExpiresAt,
		EnvironmentID: req.EnvironmentID,
	}

	for {
		err := apiKey.GenerateKey()
		if err != nil {
			return nil, err
		}

		exists, err := uc.apiKeyRepo.Exists(ctx, apiKey.Key)
		if err != nil {
			return nil, err
		}

		if !exists {
			break
		}
	}

	if err := uc.apiKeyRepo.Create(ctx, apiKey); err != nil {
		return nil, err
	}

	return &dto.APIKeyResponse{
		ID:            apiKey.ID,
		Key:           apiKey.KeySummary(),
		Status:        apiKey.Status,
		LastUsed:      apiKey.LastUsed,
		ExpiresAt:     apiKey.ExpiresAt,
		EnvironmentID: apiKey.EnvironmentID,
		CreatedAt:     apiKey.CreatedAt,
	}, nil
}

func (uc *useCase) validateReq(req *dto.APIKeyCreate) errors.Error {
	return uc.validator.ValidateStruct(
		req,
		map[string]string{
			"expires_at.utc":          "expires_at must be in UTC format",
			"environment_id.gt":       "environment_id must be greater than 0",
			"environment_id.required": "environment_id is required",
		},
	)
}

func NewUseCase(
	validator validator.Validator, apiKeyRepo APIKeyRepository,
) UseCase {
	return &useCase{
		validator:  validator,
		apiKeyRepo: apiKeyRepo,
	}
}
