package app

import (
	"context"
	"fmt"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/ports/outbound"
)

type EnvironmentUseCase struct {
	environmentRepo outbound.EnvironmentPort
	projectRepo     outbound.ProjectPort
}

func (u *EnvironmentUseCase) GetByID(
	ctx context.Context, id int,
) (*dto.EnvironmentResponse, *errors.Error) {
	environment, err := u.environmentRepo.FindByID(ctx, id)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrEnvironmentNotFound
		}
		return nil, err
	}

	serviceResp := make(
		[]*dto.EnvironmentServiceResponse, len(environment.Services),
	)
	for i, service := range environment.Services {
		serviceResp[i] = &dto.EnvironmentServiceResponse{
			ID:               service.ID,
			Name:             service.Name,
			Version:          service.Version,
			MaxRequest:       service.MaxRequest,
			AvailableRequest: service.AvailableRequest,
			AssignedAt:       service.AssignedAt,
		}
	}

	return &dto.EnvironmentResponse{
		ID:        environment.ID,
		Name:      environment.Name,
		Status:    environment.Status,
		ProjectID: environment.ProjectID,
		CreatedAt: environment.CreatedAt,
		Services:  serviceResp,
	}, nil
}

func (u *EnvironmentUseCase) AssignService(
	ctx context.Context, id int, req *dto.EnvironmentService,
) *errors.Error {
	service := entities.EnvironmentService{
		ID:               req.ID,
		MaxRequest:       req.MaxRequest,
		AvailableRequest: req.MaxRequest,
	}

	if err := service.Validate(); err != nil {
		return err
	}

	exists, err := u.environmentRepo.ExistsServiceIn(ctx, id, service.ID)
	if err != nil {
		return err
	}

	if exists {
		return errors.ErrEnvironmentServiceAlreadyExists
	}

	quota, err := u.environmentRepo.GetProjectServiceQuotaUsage(
		ctx, id, service.ID,
	)
	if err != nil {
		if err == errors.ErrNotFound {
			return errors.ErrEnvironmentNotFound
		}
		return err
	}

	if quota.MaxAllowed > -1 {
		if service.MaxRequest == -1 {
			return errors.ErrInfiniteRequestsNotAllowed
		}

		if quota.CurrentAllocated+service.MaxRequest > quota.MaxAllowed {
			return errors.ErrMaxRequestExceededForServiceInProyect
		}
	}

	return u.environmentRepo.AddService(ctx, id, &service)
}

func (u *EnvironmentUseCase) Create(
	ctx context.Context, req *dto.EnvironmentCreate,
) (*dto.EnvironmentResponse, *errors.Error) {
	exists, err := u.projectRepo.Exists(ctx, req.ProjectID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.ErrProjectNotFound
	}

	var errs []string
	services := make([]*entities.EnvironmentService, len(req.Services))
	for i, service := range req.Services {
		quota, err := u.projectRepo.GetProjectServiceQuotaUsage(
			ctx, req.ProjectID, service.ID,
		)
		if err != nil {
			if err == errors.ErrNotFound {
				err = errors.ErrServiceNotFound
			}
			err := err.AddDetail(fmt.Sprintf("Service %d", service.ID))
			errs = append(errs, err.Error())
			continue
		}

		if quota.MaxAllowed > -1 {
			if service.MaxRequest == -1 {
				err := errors.ErrInfiniteRequestsNotAllowed.
					AddDetail(fmt.Sprintf("Service %d", service.ID))
				errs = append(errs, err.Error())
				continue
			}

			if quota.CurrentAllocated+service.MaxRequest > quota.MaxAllowed {
				err := errors.ErrMaxRequestExceededForServiceInProyect.
					AddDetail(fmt.Sprintf("Service %d", service.ID))
				errs = append(errs, err.Error())
				continue
			}
		}

		services[i] = &entities.EnvironmentService{
			ID:               service.ID,
			MaxRequest:       service.MaxRequest,
			AvailableRequest: service.MaxRequest,
		}
	}

	if len(errs) > 0 {
		err := errors.NewError(
			errors.CodeValidationError,
			"Service assignment failed",
			errs...,
		)
		return nil, err
	}

	environment := entities.Environment{
		Name:      req.Name,
		Status:    enums.EnvironmentActive,
		ProjectID: req.ProjectID,
		Services:  services,
	}

	if err := environment.Validate(); err != nil {
		return nil, err
	}

	if err := u.environmentRepo.Save(ctx, &environment); err != nil {
		return nil, err
	}

	serviceResp := make(
		[]*dto.EnvironmentServiceResponse, len(environment.Services),
	)
	for i, service := range environment.Services {
		serviceResp[i] = &dto.EnvironmentServiceResponse{
			ID:               service.ID,
			Name:             service.Name,
			Version:          service.Version,
			MaxRequest:       service.MaxRequest,
			AvailableRequest: service.AvailableRequest,
			AssignedAt:       service.AssignedAt,
		}
	}

	return &dto.EnvironmentResponse{
		ID:        environment.ID,
		Name:      environment.Name,
		Status:    environment.Status,
		ProjectID: environment.ProjectID,
		CreatedAt: environment.CreatedAt,
		Services:  serviceResp,
	}, nil
}

func NewEnvironmentUseCase(
	environmentRepo outbound.EnvironmentPort,
	projectRepo outbound.ProjectPort,
) *EnvironmentUseCase {
	return &EnvironmentUseCase{
		environmentRepo: environmentRepo,
		projectRepo:     projectRepo,
	}
}
