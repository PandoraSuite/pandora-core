package app

import (
	"context"
	"errors"

	"github.com/MAD-py/pandora-core/internal/app/utils"
	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/ports/outbound"
)

type ClientUseCase struct {
	clientRepo outbound.ClientRepositoryPort
}

func (u *ClientUseCase) Create(
	ctx context.Context, req *dto.ClientCreate,
) (*dto.ClientResponse, error) {
	if req.Name == "" {
		return nil, errors.New("name of the client cannot be empty")
	}

	if utils.ValidateEmail(req.Email) {
		return nil, errors.New("invalid email")
	}

	client, err := u.clientRepo.Save(
		ctx,
		&entities.Client{
			Type:  req.Type,
			Name:  req.Name,
			Email: req.Email,
		},
	)
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
