package routes

import (
	"github.com/MAD-py/pandora-core/internal/adapters/http/bootstrap"
	"github.com/MAD-py/pandora-core/internal/adapters/http/handlers"
	"github.com/MAD-py/pandora-core/internal/app/auth"
	"github.com/gin-gonic/gin"
)

func RegisterLoginRoutes(rg *gin.RouterGroup, deps *bootstrap.Dependencies) {
	authUC := auth.NewAutenticateUseCase(
		deps.Validator, deps.TokenProvider, deps.CredentialsRepo,
	)

	auth := rg.Group("/auth")
	{
		auth.POST("/login", handlers.Authenticate(authUC))
	}
}

func RegisterAuthRoutes(rg *gin.RouterGroup, deps *bootstrap.Dependencies) {
	passChangeUC := auth.NewPasswordChangeUseCase(
		deps.Validator, deps.CredentialsRepo,
	)
	reauthenticateUC := auth.NewReauthenticateUseCase(
		deps.Validator, deps.TokenProvider, deps.CredentialsRepo,
	)

	auth := rg.Group("/auth")
	{
		auth.POST("/change-password", handlers.ChangePassword(passChangeUC))
		auth.POST("/reauthenticate", handlers.Reauthenticate(reauthenticateUC))
	}
}
