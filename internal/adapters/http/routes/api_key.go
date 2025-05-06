package routes

import (
	"github.com/MAD-py/pandora-core/internal/adapters/http/handlers"
	apikey "github.com/MAD-py/pandora-core/internal/app/api_key"
	"github.com/gin-gonic/gin"
)

func RegisterProtectedAPIKeyRoutes(rg *gin.RouterGroup) {
	createUC := apikey.NewCreateUseCase(nil, nil)
	updateUC := apikey.NewUpdateUseCase(nil, nil)

	apiKeys := rg.Group("/api-keys")
	{
		apiKeys.POST("", handlers.APIKeyCreate(createUC))
		apiKeys.PUT("/:id", handlers.APIKeyUpdate(updateUC))
	}
}
