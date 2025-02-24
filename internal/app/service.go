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

func (s *ServiceUseCase) Create(
	ctx context.Context, req *dto.ServiceCreate,
) (*dto.ServiceResponse, error) {
	if req.Name == "" {
		return nil, errors.New("name of the service cannot be empty")
	}

	service, err := s.serviceRepo.Save(
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
