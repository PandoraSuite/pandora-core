package update

import (
	"context"
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, id int, req *dto.APIKeyUpdate) (*dto.APIKeyResponse, errors.Error)
}

type useCase struct {
	validator validator.Validator

	apiKeyRepo APIKeyRepository
}

func (uc *useCase) Execute(
	ctx context.Context, id int, req *dto.APIKeyUpdate,
) (*dto.APIKeyResponse, errors.Error) {
	if err := uc.validateInput(id, req); err != nil {
		return nil, err
	}

	apiKey, err := uc.apiKeyRepo.Update(ctx, id, req)
	if err != nil {
		if err.Code() == errors.CodeNotFound {
			return nil, errors.NewEntityNotFound(
				"APIKey",
				"api key not found",
				map[string]any{"id": id},
				err,
			)
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

func (uc *useCase) validateInput(id int, req *dto.APIKeyUpdate) errors.Error {
	var err errors.Error

	if errID := uc.validateID(id); errID != nil {
		err = errors.Aggregate(err, errID)
	}

	if errReq := uc.validateReq(req); errReq != nil {
		err = errors.Aggregate(err, errReq)
	}

	return err
}

func (uc *useCase) validateID(id int) errors.Error {
	return uc.validator.ValidateVariable(
		id,
		"id",
		"required,gt=0",
		map[string]string{
			"gt":       "id must be greater than 0",
			"required": "id is required",
		},
	)
}

func (uc *useCase) validateReq(req *dto.APIKeyUpdate) errors.Error {
	var err errors.Error

	validationErr := uc.validator.ValidateStruct(
		req,
		map[string]string{
			"expires_at.utc": "expires_at must be in UTC format",
		},
	)

	if validationErr != nil {
		err = errors.Aggregate(err, validationErr)
	}

	if !req.ExpiresAt.IsZero() && req.ExpiresAt.Before(time.Now()) {
		err = errors.Aggregate(
			err,
			errors.NewAttributeValidationFailed(
				"APIKeyUpdate",
				"expires_at",
				"expires_at must be in the future",
				nil,
			),
		)
	}

	return err
}

func NewUseCase(
	validator validator.Validator, apiKeyRepo APIKeyRepository,
) UseCase {
	return &useCase{
		apiKeyRepo: apiKeyRepo,
		validator:  validator,
	}
}
