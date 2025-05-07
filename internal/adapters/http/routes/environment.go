package routes

import (
	"github.com/MAD-py/pandora-core/internal/adapters/http/handlers"
	"github.com/MAD-py/pandora-core/internal/app/environment"
	"github.com/gin-gonic/gin"
)

func RegisterEnvironmentRoutes(rg *gin.RouterGroup) {
	getUC := environment.NewGetUseCase(nil, nil)
	createUc := environment.NewCreateUseCase(nil, nil, nil)
	updateUC := environment.NewUpdateUseCase(nil, nil)
	listAPIKeysUC := environment.NewListAPIKeyUseCase(nil, nil, nil)
	resetRequestUC := environment.NewResetRequestUseCase(nil, nil)
	assignServiceUC := environment.NewAssignServiceUseCase(nil, nil)
	removeServiceUC := environment.NewRemoveServiceUseCase(nil, nil)
	updateServiceUC := environment.NewUpdateServiceUseCase(nil, nil)

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
		environments.PATCH(
			"/:id/services/:service_id/reset-requests",
			handlers.EnvironmentResetRequest(resetRequestUC),
		)
	}
}
