package routes

import (
	"github.com/MAD-py/pandora-core/internal/adapters/http/bootstrap"
	"github.com/MAD-py/pandora-core/internal/adapters/http/handlers"
	"github.com/MAD-py/pandora-core/internal/app/client"
	"github.com/gin-gonic/gin"
)

func RegisterClientRoutes(rg *gin.RouterGroup, deps *bootstrap.Dependencies) {
	getUC := client.NewGetUseCase(
		deps.Validator, deps.Repositories.Client(),
	)
	listUC := client.NewListUseCase(
		deps.Validator, deps.Repositories.Client(),
	)
	createUC := client.NewCreateUseCase(
		deps.Validator, deps.Repositories.Client(),
	)
	updateUC := client.NewUpdateUseCase(
		deps.Validator, deps.Repositories.Client(),
	)
	listProjectsUC := client.NewListProjectsUseCase(
		deps.Validator,
		deps.Repositories.Client(),
		deps.Repositories.Project(),
	)

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
