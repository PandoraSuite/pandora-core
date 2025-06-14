package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/MAD-py/pandora-core/internal/adapters/http/bootstrap"
	"github.com/MAD-py/pandora-core/internal/adapters/http/handlers"
	"github.com/MAD-py/pandora-core/internal/app/health"
)

func RegisterHealthRoutes(rg *gin.RouterGroup, deps *bootstrap.Dependencies) {
	checkUC := health.NewCheckUseCase(deps.Repositories)

	health := rg.Group("/health")
	{
		health.GET("", handlers.HealthCheck(checkUC))
	}
}
