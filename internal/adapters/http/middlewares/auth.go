package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/MAD-py/pandora-core/internal/adapters/http/errors"
	"github.com/MAD-py/pandora-core/internal/app/auth"
)

func ValidateToken(useCase auth.TokenValidationUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Header("WWW-Authenticate", `Bearer realm="Access to the API"`)
			c.Error(errors.NewUnauthorized("Authorization header missing"))
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Error(
				errors.NewUnauthorized("Invalid token type, expected 'Bearer'"),
			)
			c.Abort()
			return
		}

		username, err := useCase.Execute(c.Request.Context(), parts[0])
		if err != nil {
			c.Error(err)
			c.Abort()
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
			c.Error(errors.NewInternal("Username not found in context"))
			c.Abort()
			return
		}

		ok, err := useCase.Execute(c.Request.Context(), username)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		if ok {
			c.Error(
				errors.NewForbidden(
					"Password change required before continuing",
				),
			)
			c.Abort()
			return
		}

		c.Next()
	}
}
