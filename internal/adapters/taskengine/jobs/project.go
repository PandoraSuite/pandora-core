package jobs

import (
	"github.com/MAD-py/go-taskengine/taskengine"
	"github.com/MAD-py/pandora-core/internal/app/project"
)

func ProjectQuotaReset(useCase project.ResetDueRequestsUseCase) taskengine.Job {
	return func(ctx *taskengine.Context) error {
		ctx.Logger().Infof(
			"Starting ProjectQuotaReset job - Tick: %d", ctx.CurrentTick(),
		)

		resetResults, err := useCase.Execute(ctx)
		if err != nil {
			ctx.Logger().Errorf(
				"Error executing ProjectQuotaReset - Tick: %d - Error: %s",
				ctx.CurrentTick(), err.Error(),
			)
			return err
		}

		totalProjectsReset := len(resetResults)
		totalEnvironmentsReset := 0
		totalServicesReset := 0

		if totalProjectsReset > 0 {
			ctx.Logger().Infof(
				"Quota reset completed successfully - Projects processed: %d",
				totalProjectsReset,
			)

			for _, projectReset := range resetResults {
				environmentCount := len(projectReset.EnvironmentServices)
				totalEnvironmentsReset += environmentCount

				serviceCount := 0
				for _, envService := range projectReset.EnvironmentServices {
					if envService.Service != nil {
						serviceCount++
					}
				}
				totalServicesReset += serviceCount

				ctx.Logger().Infof(
					"Project reset - ID: %d, Name: %s, Status: %s, Environments: %d, Services: %d",
					projectReset.ID,
					projectReset.Name,
					projectReset.Status,
					environmentCount,
					serviceCount,
				)

				for _, envService := range projectReset.EnvironmentServices {
					if envService.Service != nil {
						ctx.Logger().Infof(
							"  - Environment ID: %d, Name: %s, Service: %s (ID: %d), Quota reset: %d",
							envService.ID,
							envService.Name,
							envService.Service.Name,
							envService.Service.ID,
							envService.Service.AvailableRequest,
						)
					}
				}
			}
		} else {
			ctx.Logger().Info("No projects found requiring quota reset")
		}

		ctx.Logger().Infof(
			"ProjectQuotaReset job completed - Tick: %d - Summary: %d projects, %d environments, %d services reset",
			ctx.CurrentTick(),
			totalProjectsReset,
			totalEnvironmentsReset,
			totalServicesReset,
		)

		return nil
	}
}
