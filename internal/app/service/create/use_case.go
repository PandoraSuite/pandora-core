package create

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/validator"
)

type UseCase interface {
	Execute(ctx context.Context, req *dto.ServiceCreate) (*dto.ServiceResponse, errors.Error)
}

type useCase struct {
	validator validator.Validator

	serviceRepo ServiceRepository
}

func (uc *useCase) Execute(
	ctx context.Context, req *dto.ServiceCreate,
) (*dto.ServiceResponse, errors.Error) {
	if err := uc.validateReq(req); err != nil {
		return nil, err
	}

	service := entities.Service{
		Name:    req.Name,
		Status:  enums.ServiceStatusEnabled,
		Version: req.Version,
	}

	if err := uc.serviceRepo.Create(ctx, &service); err != nil {
		if err.Code() == errors.CodeAlreadyExists {
			return nil, errors.NewEntityAlreadyExists(
				"Service",
				"Service with this name and version already exists",
				map[string]any{
					"name":    req.Name,
					"version": req.Version,
				},
				err.Unwrap(),
			)
		}
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

func (uc *useCase) validateReq(req *dto.ServiceCreate) errors.Error {
	return uc.validator.ValidateStruct(
		req,
		map[string]string{
			"name.required":    "name is required",
			"version.required": "version is required",
			"version.max":      "version must be at most 16 characters long",
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
