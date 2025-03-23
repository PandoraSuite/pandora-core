package app

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/ports/outbound"
)

type ClientUseCase struct {
	clientRepo outbound.ClientPort
}

func (c *ClientUseCase) Update(
	ctx context.Context, id int, req *dto.ClientUpdate,
) *errors.Error {
	return c.clientRepo.Update(ctx, id, req)
}

func (c *ClientUseCase) GetByID(
	ctx context.Context, id int,
) (*dto.ClientResponse, *errors.Error) {
	if id <= 0 {
		return nil, errors.ErrInvalidClientID
	}

	client, err := c.clientRepo.FindByID(ctx, id)
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

func (u *ClientUseCase) GetAll(
	ctx context.Context, req *dto.ClientFilter,
) ([]*dto.ClientResponse, *errors.Error) {
	clients, err := u.clientRepo.FindAll(ctx, req)
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

func (u *ClientUseCase) Create(
	ctx context.Context, req *dto.ClientCreate,
) (*dto.ClientResponse, *errors.Error) {
	client := entities.Client{
		Type:  req.Type,
		Name:  req.Name,
		Email: req.Email,
	}

	if err := client.Validate(); err != nil {
		return nil, err
	}

	if err := u.clientRepo.Save(ctx, &client); err != nil {
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

func NewClientUseCase(clientRepo outbound.ClientPort) *ClientUseCase {
	return &ClientUseCase{clientRepo: clientRepo}
}
