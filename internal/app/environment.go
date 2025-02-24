package app

import (
	"context"
	"errors"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/ports/outbound"
)

type EnvironmentUseCase struct {
	clientRepo outbound.EnvironmentRepositoryPort
}

func (u *EnvironmentUseCase) Create(
	ctx context.Context, req *dto.EnvironmentCreate,
) (*dto.EnvironmentResponse, error) {
	if req.Name == "" {
		return nil, errors.New("name of the environment cannot be empty")
	}

	client, err := u.clientRepo.Save(
		ctx,
		&entities.Environment{
			Name:      req.Name,
			Status:    enums.EnvironmentActive,
			ProjectID: req.ProjectID,
		},
	)
	if err != nil {
		return nil, err
	}

	return &dto.EnvironmentResponse{
		ID:        client.ID,
		Name:      client.Name,
		Status:    client.Status,
		ProjectID: client.ProjectID,
		CreatedAt: client.CreatedAt,
	}, nil
}
