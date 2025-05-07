package routes

import (
	"github.com/MAD-py/pandora-core/internal/adapters/http/handlers"
	"github.com/MAD-py/pandora-core/internal/app/service"
	"github.com/gin-gonic/gin"
)

func RegisterServiceRoutes(rg *gin.RouterGroup) {
	listUC := service.NewListUseCase(nil, nil)
	createUC := service.NewCreateUseCase(nil, nil)
	deleteUC := service.NewDeleteUseCase(nil, nil, nil, nil)
	updateStatusUC := service.NewUpdateStatusUseCase(nil, nil)

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
