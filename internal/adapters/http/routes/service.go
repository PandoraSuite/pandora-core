package routes

import (
	"github.com/MAD-py/pandora-core/internal/adapters/http/bootstrap"
	"github.com/MAD-py/pandora-core/internal/adapters/http/handlers"
	"github.com/MAD-py/pandora-core/internal/app/service"
	"github.com/gin-gonic/gin"
)

func RegisterServiceRoutes(rg *gin.RouterGroup, deps *bootstrap.Dependencies) {
	listUC := service.NewListUseCase(
		deps.Validator, deps.Repositories.Service(),
	)
	createUC := service.NewCreateUseCase(
		deps.Validator, deps.Repositories.Service(),
	)
	deleteUC := service.NewDeleteUseCase(
		deps.Validator,
		deps.Repositories.Service(),
		deps.Repositories.Project(),
	)
	updateStatusUC := service.NewUpdateStatusUseCase(
		deps.Validator, deps.Repositories.Service(),
	)

	services := rg.Group("/services")
	{
		services.GET("", handlers.ServiceList(listUC))
		services.POST("", handlers.ServiceCreate(createUC))
		services.DELETE("/:id", handlers.ServiceDelete(deleteUC))
		services.PATCH(
			"/:id/status",
			handlers.ServiceUpdateStatus(updateStatusUC),
		)
	}
}
