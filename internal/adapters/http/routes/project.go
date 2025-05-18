package routes

import (
	"github.com/MAD-py/pandora-core/internal/adapters/http/bootstrap"
	"github.com/MAD-py/pandora-core/internal/adapters/http/handlers"
	"github.com/MAD-py/pandora-core/internal/app/project"
	"github.com/gin-gonic/gin"
)

func RegisterProjectRoutes(rg *gin.RouterGroup, deps *bootstrap.Dependencies) {
	getUC := project.NewGetUseCase(
		deps.Validator, deps.Repositories.Project(),
	)
	listUC := project.NewListUseCase(deps.Repositories.Project())
	createUC := project.NewCreateUseCase(
		deps.Validator, deps.Repositories.Project(),
	)
	updateUC := project.NewUpdateUseCase(
		deps.Validator, deps.Repositories.Project(),
	)
	resetRequest := project.NewResetRequestUseCase(
		deps.Validator, deps.Repositories.Project(),
	)
	assignService := project.NewAssignServiceUseCase(
		deps.Validator, deps.Repositories.Project(),
	)
	removeService := project.NewRemoveServiceUseCase(
		deps.Validator,
		deps.Repositories.Project(),
		deps.Repositories.Environment(),
	)
	updateService := project.NewUpdateServiceUseCase(
		deps.Validator,
		deps.Repositories.Project(),
		deps.Repositories.Environment(),
	)
	listEvntiromentsUC := project.NewListEnvironmentsUseCase(
		deps.Validator,
		deps.Repositories.Project(),
		deps.Repositories.Environment(),
	)

	projects := rg.Group("/projects")
	{
		projects.GET("", handlers.ProjectList(listUC))
		projects.POST("", handlers.ProjectCreate(createUC))
		projects.GET("/:id", handlers.ProjectGet(getUC))
		projects.PATCH("/:id", handlers.ProjectUpdate(updateUC))
		projects.GET(
			"/:id/environments",
			handlers.ProjectListEnvironments(listEvntiromentsUC),
		)
		projects.POST(
			"/:id/services",
			handlers.ProjectAssignService(assignService),
		)
		projects.DELETE(
			"/:id/services/:service_id",
			handlers.ProjectRemoveService(removeService),
		)
		projects.PATCH(
			"/:id/services/:service_id",
			handlers.ProjectUpdateService(updateService),
		)
		projects.POST(
			"/:id/services/:service_id/reset-requests",
			handlers.ProjectResetRequest(resetRequest),
		)
	}
}
