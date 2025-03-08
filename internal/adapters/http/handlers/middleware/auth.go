package middleware

import (
	"net/http"
	"strings"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/ports/inbound"
	"github.com/gin-gonic/gin"
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
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("username", username)
		c.Next()
	}
}
