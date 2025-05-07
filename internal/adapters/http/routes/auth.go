package routes

import (
	"github.com/MAD-py/pandora-core/internal/adapters/http/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterLoginRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", handlers.Authenticate(nil))
	}
}

func RegisterAuthRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/change-password", handlers.ChangePassword(nil))
	}
}
