package app

import (
	"context"
	"errors"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/ports/outbound"
)

type ProjectUseCase struct {
	clientRepo outbound.ProjectRepositoryPort
}

func (u *ProjectUseCase) Create(
	ctx context.Context, req *dto.ProjectCreate,
) (*dto.ProjectResponse, error) {
	if req.Name == "" {
		return nil, errors.New("name of the project cannot be empty")
	}

	client, err := u.clientRepo.Save(
		ctx,
		&entities.Project{
			Name:     req.Name,
			Status:   req.Status,
			ClientID: req.ClientID,
		},
	)
	if err != nil {
		return nil, err
	}

	return &dto.ProjectResponse{
		ID:        client.ID,
		Name:      client.Name,
		Status:    client.Status,
		ClientID:  client.ClientID,
		CreatedAt: client.CreatedAt,
	}, nil
}
