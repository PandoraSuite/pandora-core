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

func (u *EnvironmentUseCase) UpdateService(
	ctx context.Context, id, serviceID int, req *dto.EnvironmentServiceUpdate,
) (*dto.EnvironmentServiceResponse, *errors.Error) {
	if req.MaxRequest < -1 {
		return nil, errors.ErrInvalidMaxRequest
	}

	exists, err := u.environmentRepo.Exists(ctx, id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.ErrEnvironmentNotFound
	}

	service, err := u.environmentRepo.FindServiceByID(
		ctx, id, serviceID,
	)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrEnvironmentServiceNotFound
		}
		return nil, err
	}

	quota, err := u.environmentRepo.GetProjectServiceQuotaUsage(
		ctx, id, serviceID,
	)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrServiceNotAssignedToProject
		}
		return nil, err
	}

	if req.MaxRequest == -1 {
		if quota.MaxAllowed != -1 {
			return nil, errors.ErrInfiniteRequestsNotAllowed
		}
	} else if quota.MaxAllowed != -1 {
		newAllocated := quota.CurrentAllocated - service.MaxRequest + req.MaxRequest
		if newAllocated > quota.MaxAllowed {
			return nil, errors.ErrMaxRequestExceededForServiceInProyect
		}
	}

	service, err = u.environmentRepo.UpdateService(
		ctx, id, serviceID, req,
	)
	if err != nil {
		return nil, err
	}

	return &dto.EnvironmentServiceResponse{
		ID:               service.ID,
		Name:             service.Name,
		Version:          service.Version,
		MaxRequest:       service.MaxRequest,
		AvailableRequest: service.AvailableRequest,
		AssignedAt:       service.AssignedAt,
	}, nil
}

func (u *EnvironmentUseCase) Update(
	ctx context.Context, id int, req *dto.EnvironmentUpdate,
) (*dto.EnvironmentResponse, *errors.Error) {
	environment, err := u.environmentRepo.Update(ctx, id, req)
	if err != nil {
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

func (u *EnvironmentUseCase) ResetServiceRequests(
	ctx context.Context, id, serviceID int,
) (*dto.EnvironmentServiceResponse, *errors.Error) {
	exists, err := u.environmentRepo.Exists(ctx, id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.ErrEnvironmentNotFound
	}

	service, err := u.environmentRepo.ResetAvailableRequests(ctx, id, serviceID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrServiceNotAssignedToEnvironment
		}
		return nil, err
	}

	return &dto.EnvironmentServiceResponse{
		ID:               service.ID,
		Name:             service.Name,
		Version:          service.Version,
		MaxRequest:       service.MaxRequest,
		AvailableRequest: service.AvailableRequest,
		AssignedAt:       service.AssignedAt,
	}, nil
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

func (u *EnvironmentUseCase) RemoveService(
	ctx context.Context, id, serviceID int,
) *errors.Error {
	exists, err := u.environmentRepo.Exists(ctx, id)
	if err != nil {
		return err
	}

	if !exists {
		return errors.ErrEnvironmentNotFound
	}

	n, err := u.environmentRepo.RemoveService(ctx, id, serviceID)
	if err != nil {
		return err
	}

	if n == 0 {
		return errors.ErrServiceNotFound
	}

	return nil
}

func (u *EnvironmentUseCase) AssignService(
	ctx context.Context, id int, req *dto.EnvironmentService,
) (*dto.EnvironmentServiceResponse, *errors.Error) {
	service := entities.EnvironmentService{
		ID:               req.ID,
		MaxRequest:       req.MaxRequest,
		AvailableRequest: req.MaxRequest,
	}

	if err := service.Validate(); err != nil {
		return nil, err
	}

	exists, err := u.environmentRepo.Exists(ctx, id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.ErrEnvironmentNotFound
	}

	exists, err = u.environmentRepo.ExistsServiceIn(ctx, id, service.ID)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.ErrEnvironmentServiceAlreadyExists
	}

	quota, err := u.environmentRepo.GetProjectServiceQuotaUsage(
		ctx, id, service.ID,
	)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrServiceNotAssignedToProject
		}
		return nil, err
	}

	if quota.MaxAllowed > -1 {
		if service.MaxRequest == -1 {
			return nil, errors.ErrInfiniteRequestsNotAllowed
		}

		if quota.CurrentAllocated+service.MaxRequest > quota.MaxAllowed {
			return nil, errors.ErrMaxRequestExceededForServiceInProyect
		}
	}

	if err := u.environmentRepo.AddService(ctx, id, &service); err != nil {
		return nil, err
	}

	return &dto.EnvironmentServiceResponse{
		ID:               service.ID,
		Name:             service.Name,
		Version:          service.Version,
		MaxRequest:       service.MaxRequest,
		AvailableRequest: service.AvailableRequest,
		AssignedAt:       service.AssignedAt,
	}, nil
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
