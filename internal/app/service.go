package app

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/ports/outbound"
)

type ServiceUseCase struct {
	serviceRepo outbound.ServicePort
}

func (u *ServiceUseCase) GetServices(ctx context.Context) ([]*dto.ServiceResponse, *errors.Error) {
	services, err := u.serviceRepo.FindAll(ctx)
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

func (u *ServiceUseCase) GetActiveServices(ctx context.Context) ([]*dto.ServiceResponse, *errors.Error) {
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
) (*dto.ServiceResponse, *errors.Error) {
	service := entities.Service{
		Name:    req.Name,
		Status:  enums.ServiceActive,
		Version: req.Version,
	}

	if err := service.Validate(); err != nil {
		return nil, err
	}

	if err := u.serviceRepo.Save(ctx, &service); err != nil {
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

func NewServiceUseCase(serviceRepo outbound.ServicePort) *ServiceUseCase {
	return &ServiceUseCase{serviceRepo: serviceRepo}
}
