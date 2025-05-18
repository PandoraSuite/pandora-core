package update

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, id int, req *dto.ClientUpdate) (*dto.ClientResponse, errors.Error)
}

type useCase struct {
	validator validator.Validator

	clientRepo ClientRepository
}

func (uc *useCase) Execute(ctx context.Context, id int, req *dto.ClientUpdate) (*dto.ClientResponse, errors.Error) {
	if err := uc.validateReq(req); err != nil {
		return nil, err
	}

	client, err := uc.clientRepo.Update(ctx, id, req)
	if err != nil {
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

func (uc *useCase) validateReq(req *dto.ClientUpdate) errors.Error {
	return uc.validator.ValidateStruct(
		req,
		map[string]string{
			"type.enums":  "type must be one of the following: developer, organization",
			"email.email": "email must be a valid email address",
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
