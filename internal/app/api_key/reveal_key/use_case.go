package revealkey

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, id int) (*dto.APIKeyRevealKeyResponse, errors.Error)
}

type useCase struct {
	validator validator.Validator

	apiKeyRepo APIKeyRepository
}

func (uc *useCase) Execute(ctx context.Context, id int) (*dto.APIKeyRevealKeyResponse, errors.Error) {
	if err := uc.validateID(id); err != nil {
		return nil, err
	}

	apiKey, err := uc.apiKeyRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.APIKeyRevealKeyResponse{
		Key: apiKey.Key,
	}, nil
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

func NewUseCase(
	validator validator.Validator, apiKeyRepo APIKeyRepository,
) UseCase {
	return &useCase{
		validator:  validator,
		apiKeyRepo: apiKeyRepo,
	}
}
