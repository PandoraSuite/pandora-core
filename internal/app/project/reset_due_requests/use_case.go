package resetduerequests

import (
	"context"
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/utils"
)

type UseCase interface {
	Execute(ctx context.Context) ([]*dto.ProjectReset, errors.Error)
}

type useCase struct {
	projectRepo ProjectRepository
}

func (uc *useCase) Execute(ctx context.Context) ([]*dto.ProjectReset, errors.Error) {
	today := utils.TruncateToDay(time.Now())
	projects, err := uc.projectRepo.ListProjectServiceDueForReset(
		ctx, today,
	)
	if err != nil {
		return nil, err
	}

	projectsResponse := make([]*dto.ProjectReset, 0)

	var errs errors.Error
	for _, project := range projects {
		project.CalculateNextServicesReset()

		envServices := make([]*dto.EnvironmentServiceReset, 0)
		for _, service := range project.Services {
			resp, err := uc.projectRepo.ResetProjectServiceUsage(
				ctx, project.ID, service.ID, service.NextReset,
			)
			if err != nil {
				errs = errors.Aggregate(errs, err)
				continue
			}

			envServices = append(envServices, resp...)
		}

		projectsResponse = append(projectsResponse, &dto.ProjectReset{
			ID:                  project.ID,
			Name:                project.Name,
			Status:              project.Status,
			EnvironmentServices: envServices,
		})
	}

	return projectsResponse, nil
}

func NewUseCase(projectRepo ProjectRepository) UseCase {
	return &useCase{
		projectRepo: projectRepo,
	}
}
