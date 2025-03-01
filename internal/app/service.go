package app

import (
	"context"
	"errors"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/ports/outbound"
)

type ServiceUseCase struct {
	serviceRepo outbound.ServiceRepositoryPort
}

func (u *ServiceUseCase) GetActiveServices(ctx context.Context) ([]*dto.ServiceResponse, error) {
	services, err := u.serviceRepo.FindActiveServices(ctx)
	if err != nil {
		return nil, err
	}

	serviceResponses := make([]*dto.ServiceResponse, len(services))
	for i, service := range services {
		serviceResponses[i] = &dto.ServiceResponse{
			ID:        service.ID,
			Name:      service.Name,
			Status:    service.Status,
			Version:   service.Version,
			CreatedAt: service.CreatedAt,
		}
	}

	return serviceResponses, nil
}

func (u *ServiceUseCase) Create(
	ctx context.Context, req *dto.ServiceCreate,
) (*dto.ServiceResponse, error) {
	if req.Name == "" {
		return nil, errors.New("name of the service cannot be empty")
	}

	service, err := u.serviceRepo.Save(
		ctx,
		&entities.Service{
			Name:    req.Name,
			Status:  enums.ServiceActive,
			Version: req.Version,
		},
	)
	if err != nil {
		return nil, err
	}

	return &dto.ServiceResponse{
		ID:        service.ID,
		Name:      service.Name,
		Status:    service.Status,
		Version:   service.Version,
		CreatedAt: service.CreatedAt,
	}, nil
}
