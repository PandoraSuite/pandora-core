package list

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, req *dto.ClientFilter) ([]*dto.ClientResponse, errors.Error)
}

type useCase struct {
	validator validator.Validator

	clientRepo ClientRepository
}

func (uc *useCase) Execute(
	ctx context.Context, req *dto.ClientFilter,
) ([]*dto.ClientResponse, errors.Error) {
	if err := uc.validateReq(req); err != nil {
		return nil, err
	}

	clients, err := uc.clientRepo.List(ctx, req)
	if err != nil {
		return nil, err
	}

	clientResponses := make([]*dto.ClientResponse, len(clients))
	for i, client := range clients {
		clientResponses[i] = &dto.ClientResponse{
			ID:        client.ID,
			Type:      client.Type,
			Name:      client.Name,
			Email:     client.Email,
			CreatedAt: client.CreatedAt,
		}
	}

	return clientResponses, nil
}

func (uc *useCase) validateReq(req *dto.ClientFilter) errors.Error {
	return uc.validator.ValidateStruct(
		req,
		map[string]string{
			"type.enums": "type must be one of the following: developer, organization",
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
