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

	maxRequest, err := u.environmentRepo.
		GetMaxRequestByEnvironmentAndServiceInProject(ctx, id, service.ID)
	if err != nil {
		return err
	}

	if maxRequest > 0 {
		maxRequests, err := u.environmentRepo.
			ListMaxRequestsByEnvironmentAndService(ctx, id, service.ID)
		if err != nil {
			return err
		}

		var totalMaxRequest int
		for _, v := range maxRequests {
			totalMaxRequest += v
		}

		if service.MaxRequest+totalMaxRequest > maxRequest {
			return errors.ErrMaxRequestExceededForServiceInProyect
		}
	}

	return u.environmentRepo.AddService(ctx, id, &service)
}

func (u *EnvironmentUseCase) Create(
	ctx context.Context, req *dto.EnvironmentCreate,
) (*dto.EnvironmentResponse, *errors.Error) {
	project, err := u.projectRepo.FindByID(ctx, req.ProjectID)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, errors.ErrProjectNotFound
		}
		return nil, err
	}

	var errs []string
	services := make([]*entities.EnvironmentService, len(req.Services))
	for i, service := range req.Services {
		projectService, err := project.GetService(service.ID)
		if err != nil {
			err := err.AddDetail(fmt.Sprintf("Service %d", service.ID))
			errs = append(errs, err.Error())
			services[i] = nil
			continue
		}

		if projectService.MaxRequest > 0 {
			if service.MaxRequest == 0 {
				err := errors.ErrInfiniteRequestsNotAllowed.
					AddDetail(fmt.Sprintf("Service %d", service.ID))
				errs = append(errs, err.Error())
				services[i] = nil
				continue
			}

			maxRequests, err := u.environmentRepo.
				ListMaxRequestsByProjectAndService(
					ctx, req.ProjectID, service.ID,
				)
			if err != nil {
				err := err.AddDetail(fmt.Sprintf("Service %d", service.ID))
				errs = append(errs, err.Error())
				services[i] = nil
				continue
			}

			var totalMaxRequest int
			for _, v := range maxRequests {
				totalMaxRequest += v
			}

			if service.MaxRequest+totalMaxRequest > projectService.MaxRequest {
				err := errors.ErrMaxRequestExceededForServiceInProyect.
					AddDetail(fmt.Sprintf("Service %d", service.ID))
				errs = append(errs, err.Error())
				services[i] = nil
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
			ID:         service.ID,
			Name:       service.Name,
			Version:    service.Version,
			MaxRequest: service.MaxRequest,
			AssignedAt: service.AssignedAt,
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
