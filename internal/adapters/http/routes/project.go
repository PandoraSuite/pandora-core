package routes

import (
	"github.com/MAD-py/pandora-core/internal/adapters/http/handlers"
	"github.com/MAD-py/pandora-core/internal/app/project"
	"github.com/gin-gonic/gin"
)

func RegisterProjectRoutes(rg *gin.RouterGroup) {
	getUC := project.NewGetUseCase(nil, nil)
	listUC := project.NewListUseCase(nil)
	createUC := project.NewCreateUseCase(nil, nil)
	updateUC := project.NewUpdateUseCase(nil, nil)
	resetRequest := project.NewResetRequestUseCase(nil, nil)
	assignService := project.NewAssignServiceUseCase(nil, nil)
	removeService := project.NewRemoveServiceUseCase(nil, nil, nil)
	updateService := project.NewUpdateServiceUseCase(nil, nil, nil)
	listEvntiromentsUC := project.NewListEnvironmentsUseCase(nil, nil, nil)

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
