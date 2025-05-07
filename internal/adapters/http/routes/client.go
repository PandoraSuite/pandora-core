package routes

import (
	"github.com/MAD-py/pandora-core/internal/adapters/http/handlers"
	"github.com/MAD-py/pandora-core/internal/app/client"
	"github.com/gin-gonic/gin"
)

func RegisterClientRoutes(rg *gin.RouterGroup) {
	getUC := client.NewGetUseCase(nil, nil)
	listUC := client.NewListUseCase(nil, nil)
	createUC := client.NewCreateUseCase(nil, nil)
	updateUC := client.NewUpdateUseCase(nil, nil)
	listProjectsUC := client.NewListProjectsUseCase(nil, nil, nil)

	clients := rg.Group("/clients")
	{
		clients.GET("", handlers.ClientList(listUC))
		clients.POST("", handlers.ClientCreate(createUC))
		clients.GET("/:id", handlers.ClientGet(getUC))
		clients.PATCH("/:id", handlers.ClientUpdate(updateUC))
		clients.GET(
			"/:id/projects",
			handlers.ClientListProjects(listProjectsUC),
		)
	}
}
