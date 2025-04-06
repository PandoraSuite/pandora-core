package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/MAD-py/pandora-core/internal/adapters/http/handlers/utils"
	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/ports/inbound"
)

func ValidateToken(authService inbound.AuthHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Header("WWW-Authenticate", `Bearer realm="Access to the API"`)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header missing"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token type, expected 'Bearer'"})
			c.Abort()
			return
		}

		req := &dto.TokenRequest{
			Key:  parts[1],
			Type: parts[0],
		}

		username, err := authService.ValidateToken(c.Request.Context(), req)
		if err != nil {
			c.JSON(
				utils.GetDomainErrorStatusCode(err),
				utils.ErrorResponse{Error: err},
			)
			c.Abort()
			return
		}

		c.Set("username", username)
		c.Next()
	}
}

func ForcePasswordReset(authService inbound.AuthHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetString("username")
		if username == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "username not found in context"})
			c.Abort()
			return
		}

		ok, err := authService.IsPasswordResetRequired(
			c.Request.Context(), username,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if ok {
			c.JSON(
				http.StatusForbidden,
				gin.H{"error": "Password change required before continuing"},
			)
			c.Abort()
			return
		}

		c.Next()
	}
}
