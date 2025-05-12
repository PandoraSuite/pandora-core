package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/MAD-py/pandora-core/internal/adapters/http/dto"
	"github.com/MAD-py/pandora-core/internal/adapters/http/handlers/utils"
	"github.com/MAD-py/pandora-core/internal/app/auth"
)

func ValidateToken(useCase auth.TokenValidationUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Header("WWW-Authenticate", `Bearer realm="Access to the API"`)
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "authorization header missing"},
			)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "invalid token type, expected 'Bearer'"},
			)
			return
		}

		req := dto.TokenValidation{
			TokenType:   parts[0],
			AccessToken: parts[1],
		}

		username, err := useCase.Execute(c.Request.Context(), req.ToDomain())
		if err != nil {
			c.AbortWithStatusJSON(
				utils.GetDomainErrorStatusCode(err),
				utils.ErrorResponse{Error: err},
			)
			return
		}

		c.Set("username", username)
		c.Next()
	}
}

func ForcePasswordReset(useCase auth.ResetPasswordUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetString("username")
		if username == "" {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"error": "username not found in context"},
			)
			return
		}

		ok, err := useCase.Execute(c.Request.Context(), username)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"error": err.Error()},
			)
			return
		}

		if ok {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{"error": "Password change required before continuing"},
			)
			return
		}

		c.Next()
	}
}
