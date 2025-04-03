package app

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/ports/outbound"
)

type ClientUseCase struct {
	clientRepo  outbound.ClientPort
	projectRepo outbound.ProjectPort
}

func (c *ClientUseCase) Update(
	ctx context.Context, id int, req *dto.ClientUpdate,
) *errors.Error {
	return c.clientRepo.Update(ctx, id, req)
}

func (c *ClientUseCase) GetByID(
	ctx context.Context, id int,
) (*dto.ClientResponse, *errors.Error) {
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

func (u *ClientUseCase) GetProjects(
	ctx context.Context, id int,
) ([]*dto.ProjectResponse, *errors.Error) {
	exists, err := u.clientRepo.Exists(ctx, id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.ErrClientNotFound
	}

	projects, err := u.projectRepo.FindByClient(ctx, id)
	if err != nil {
		return nil, err
	}

	projectResponses := make([]*dto.ProjectResponse, len(projects))
	for i, project := range projects {
		serviceResp := make(
			[]*dto.ProjectServiceResponse, len(project.Services),
		)
		for i, service := range project.Services {
			serviceResp[i] = &dto.ProjectServiceResponse{
				ID:             service.ID,
				Name:           service.Name,
				Version:        service.Version,
				NextReset:      service.NextReset,
				MaxRequest:     service.MaxRequest,
				ResetFrequency: service.ResetFrequency,
				AssignedAt:     service.AssignedAt,
			}
		}

		projectResponses[i] = &dto.ProjectResponse{
			ID:        project.ID,
			Name:      project.Name,
			Status:    project.Status,
			ClientID:  project.ClientID,
			CreatedAt: project.CreatedAt,
			Services:  serviceResp,
		}
	}

	return projectResponses, nil
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

func NewClientUseCase(
	clientRepo outbound.ClientPort,
	projectRepo outbound.ProjectPort,
) *ClientUseCase {
	return &ClientUseCase{
		clientRepo:  clientRepo,
		projectRepo: projectRepo,
	}
}
