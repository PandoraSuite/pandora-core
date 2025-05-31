package routes

import (
	"github.com/MAD-py/pandora-core/internal/adapters/http/bootstrap"
	"github.com/MAD-py/pandora-core/internal/adapters/http/handlers"
	"github.com/MAD-py/pandora-core/internal/app/environment"
	"github.com/gin-gonic/gin"
)

func RegisterEnvironmentRoutes(rg *gin.RouterGroup, deps *bootstrap.Dependencies) {
	getUC := environment.NewGetUseCase(
		deps.Validator, deps.Repositories.Environment(),
	)
	createUc := environment.NewCreateUseCase(
		deps.Validator,
		deps.Repositories.Project(),
		deps.Repositories.Environment(),
	)
	updateUC := environment.NewUpdateUseCase(
		deps.Validator, deps.Repositories.Environment(),
	)
	listAPIKeysUC := environment.NewListAPIKeyUseCase(
		deps.Validator,
		deps.Repositories.APIKey(),
		deps.Repositories.Environment(),
	)
	resetRequestUC := environment.NewResetRequestUseCase(
		deps.Validator, deps.Repositories.Environment(),
	)
	assignServiceUC := environment.NewAssignServiceUseCase(
		deps.Validator, deps.Repositories.Environment(),
	)
	removeServiceUC := environment.NewRemoveServiceUseCase(
		deps.Validator, deps.Repositories.Environment(),
	)
	updateServiceUC := environment.NewUpdateServiceUseCase(
		deps.Validator, deps.Repositories.Environment(),
	)

	environments := rg.Group("/environments")
	{
		environments.POST(
			"", handlers.EnvironmentCreate(createUc),
		)
		environments.GET(
			"/:id", handlers.EnvironmentGet(getUC),
		)
		environments.PATCH(
			"/:id", handlers.EnvironmentUpdate(updateUC),
		)
		environments.GET(
			"/:id/api-keys",
			handlers.EnvironmentListAPIKeys(listAPIKeysUC),
		)
		environments.POST(
			"/:id/services",
			handlers.EnvironmentAssignService(assignServiceUC),
		)
		environments.DELETE(
			"/:id/services/:service_id",
			handlers.EnvironmentRemoveService(removeServiceUC),
		)
		environments.PATCH(
			"/:id/services/:service_id",
			handlers.EnvironmentUpdateService(updateServiceUC),
		)
		environments.POST(
			"/:id/services/:service_id/reset-requests",
			handlers.EnvironmentResetRequest(resetRequestUC),
		)
	}
}
