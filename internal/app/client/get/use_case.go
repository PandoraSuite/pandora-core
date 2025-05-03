package get

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, id int) (*dto.ClientResponse, errors.Error)
}

type useCase struct {
	validator validator.Validator

	clientRepo ClientRepository
}

func (uc *useCase) Execute(
	ctx context.Context, id int,
) (*dto.ClientResponse, errors.Error) {
	if err := uc.validateID(id); err != nil {
		return nil, err
	}

	client, err := uc.clientRepo.GetByID(ctx, id)
	if err != nil {
		if err.Code() == errors.CodeNotFound {
			return nil, errors.NewEntityNotFound(
				"client",
				"client not found",
				map[string]any{"id": id},
			)
		}
		return nil, err
	}

	return &dto.ClientResponse{
		ID:        client.ID,
		Type:      client.Type,
		Name:      client.Name,
		Email:     client.Email,
		CreatedAt: client.CreatedAt,
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
	validator validator.Validator, clientRepo ClientRepository,
) UseCase {
	return &useCase{
		validator:  validator,
		clientRepo: clientRepo,
	}
}
