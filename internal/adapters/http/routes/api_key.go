package routes

import (
	"github.com/MAD-py/pandora-core/internal/adapters/http/bootstrap"
	"github.com/MAD-py/pandora-core/internal/adapters/http/handlers"
	apikey "github.com/MAD-py/pandora-core/internal/app/api_key"
	"github.com/gin-gonic/gin"
)

func RegisterAPIKeyRoutes(rg *gin.RouterGroup, deps *bootstrap.Dependencies) {
	createUC := apikey.NewCreateUseCase(
		deps.Validator, deps.Repositories.APIKey(),
	)
	updateUC := apikey.NewUpdateUseCase(
		deps.Validator, deps.Repositories.APIKey(),
	)

	apiKeys := rg.Group("/api-keys")
	{
		apiKeys.POST("", handlers.APIKeyCreate(createUC))
		apiKeys.PUT("/:id", handlers.APIKeyUpdate(updateUC))
	}
}
