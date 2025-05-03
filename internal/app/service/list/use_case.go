package list

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, req *dto.ServiceFilter) ([]*dto.ServiceResponse, errors.Error)
}

type useCase struct {
	validator validator.Validator

	serviceRepo ServiceRepository
}

func (uc *useCase) Execute(
	ctx context.Context, req *dto.ServiceFilter,
) ([]*dto.ServiceResponse, errors.Error) {
	if err := uc.validateReq(req); err != nil {
		return nil, err
	}

	services, err := uc.serviceRepo.List(ctx, req)
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

func (uc *useCase) validateReq(req *dto.ServiceFilter) errors.Error {
	return uc.validator.ValidateStruct(
		req,
		map[string]string{
			"status.enums": "status must be one of the following: active, deactivated, deprecated",
		},
	)
}

func NewUseCase(
	validator validator.Validator, serviceRepo ServiceRepository,
) UseCase {
	return &useCase{
		validator:   validator,
		serviceRepo: serviceRepo,
	}
}
