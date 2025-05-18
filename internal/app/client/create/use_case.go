package create

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, req *dto.ClientCreate) (*dto.ClientResponse, errors.Error)
}

type useCase struct {
	validator validator.Validator

	clientRepo ClientRepository
}

func (uc *useCase) Execute(
	ctx context.Context, req *dto.ClientCreate,
) (*dto.ClientResponse, errors.Error) {
	if err := uc.validateReq(req); err != nil {
		return nil, err
	}

	client := entities.Client{
		Type:  req.Type,
		Name:  req.Name,
		Email: req.Email,
	}

	if err := uc.clientRepo.Create(ctx, &client); err != nil {
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

func (uc *useCase) validateReq(req *dto.ClientCreate) errors.Error {
	return uc.validator.ValidateStruct(
		req,
		map[string]string{
			"type.enums":     "type must be one of the following: developer, organization",
			"email.email":    "email must be a valid email address",
			"type.required":  "type is required",
			"name.required":  "name is required",
			"email.required": "email is required",
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
