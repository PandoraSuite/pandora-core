package routes

import (
	"github.com/MAD-py/pandora-core/internal/adapters/http/bootstrap"
	"github.com/MAD-py/pandora-core/internal/adapters/http/handlers"
	"github.com/MAD-py/pandora-core/internal/adapters/http/middlewares"
	apikey "github.com/MAD-py/pandora-core/internal/app/api_key"
	"github.com/MAD-py/pandora-core/internal/app/auth"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
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

func RegisterAPIKeySensitiveRoutes(
	rg *gin.RouterGroup, deps *bootstrap.Dependencies, middleware ...gin.HandlerFunc,
) {
	revealKeyUC := apikey.NewRevealKeyUseCase(
		deps.Validator, deps.Repositories.APIKey(),
	)

	apiKey := rg.Group("/api-keys")
	{

		{
			revealKeyHandlers := []gin.HandlerFunc{
				middlewares.ValidateScopedToken(
					auth.NewScopedTokenValidationUseCase(
						deps.Validator, deps.TokenProvider,
					),
					enums.ScopeRevealAPIKey,
				),
			}
			revealKeyHandlers = append(revealKeyHandlers, middleware...)
			revealKeyHandlers = append(revealKeyHandlers, handlers.APIKeyRevealKey(revealKeyUC))
			apiKey.GET("/:id/reveal/key", revealKeyHandlers...)
		}
	}
}
