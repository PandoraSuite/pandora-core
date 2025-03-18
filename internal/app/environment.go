package app

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/ports/outbound"
)

type EnvironmentUseCase struct {
	environmentRepo        outbound.EnvironmentPort
	projectServiceRepo     outbound.ProjectServiceFindPort
	environmentServiceRepo outbound.EnvironmentServicePort
}

func (u *EnvironmentUseCase) AssignService(
	ctx context.Context, req *dto.AssignServiceToEnvironment,
) error {
	environment, err := u.environmentRepo.FindByID(ctx, req.EnvironmentID)
	if err != nil {
		if err == errors.ErrNotFound {
			err = errors.ErrEnvironmentNotFound
		}
		return err
	}

	projectService, err := u.projectServiceRepo.FindByProjectAndService(
		ctx, environment.ProjectID, req.ServiceID,
	)
	if err != nil {
		if err == errors.ErrNotFound {
			err = errors.ErrProjectServiceNotFound
		}
		return err
	}

	if projectService.MaxRequest > 0 {
		environmentServices, err := u.environmentServiceRepo.FindByProjectAndService(
			ctx, environment.ProjectID, req.ServiceID,
		)
		if err != nil {
			if err == errors.ErrNotFound {
				err = errors.ErrEnvironmentServiceNotFound
			}
			return err
		}

		var totalMaxRequest int
		for _, s := range environmentServices {
			totalMaxRequest += s.MaxRequest
		}

		if req.MaxRequest+totalMaxRequest > projectService.MaxRequest {
			return errors.ErrMaxRequestExceededForServiceInProyect
		}
	}

	err = u.environmentServiceRepo.Save(
		ctx,
		&entities.EnvironmentService{
			ServiceID:     req.ServiceID,
			EnvironmentID: req.EnvironmentID,
			MaxRequest:    req.MaxRequest,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (u *EnvironmentUseCase) GetByProject(
	ctx context.Context, projectID int,
) ([]*dto.EnvironmentResponse, error) {
	environments, err := u.environmentRepo.FindByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	environmentResponses := make([]*dto.EnvironmentResponse, len(environments))
	for i, environment := range environments {
		environmentResponses[i] = &dto.EnvironmentResponse{
			ID:        environment.ID,
			Name:      environment.Name,
			Status:    environment.Status,
			ProjectID: environment.ProjectID,
			CreatedAt: environment.CreatedAt,
		}
	}

	return environmentResponses, nil
}

func (u *EnvironmentUseCase) Create(
	ctx context.Context, req *dto.EnvironmentCreate,
) (*dto.EnvironmentResponse, error) {
	environment := entities.Environment{
		Name:      req.Name,
		Status:    enums.EnvironmentActive,
		ProjectID: req.ProjectID,
	}

	if err := environment.Validate(); err != nil {
		return nil, err
	}

	if err := u.environmentRepo.Save(ctx, &environment); err != nil {
		return nil, err
	}

	return &dto.EnvironmentResponse{
		ID:        environment.ID,
		Name:      environment.Name,
		Status:    environment.Status,
		ProjectID: environment.ProjectID,
		CreatedAt: environment.CreatedAt,
	}, nil
}

func NewEnvironmentUseCase(
	environmentRepo outbound.EnvironmentPort,
	projectServiceRepo outbound.ProjectServiceFindPort,
	environmentServiceRepo outbound.EnvironmentServicePort,
) *EnvironmentUseCase {
	return &EnvironmentUseCase{
		environmentRepo:        environmentRepo,
		projectServiceRepo:     projectServiceRepo,
		environmentServiceRepo: environmentServiceRepo,
	}
}
