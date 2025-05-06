package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/MAD-py/pandora-core/internal/adapters/http/dto"
	"github.com/MAD-py/pandora-core/internal/adapters/http/handlers/utils"
	"github.com/MAD-py/pandora-core/internal/ports/inbound"
)

// ChangePassword godoc
// @Summary Change password
// @Description Allows an authenticated user to change their password.
// @Tags Authentication
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param request body dto.ChangePassword true "New password and confirmation"
// @Success 204
// @Failure default {object} utils.ErrorResponse "Default error response for all failures"
// @Router /api/v1/auth/change-password [post]
func ChangePassword(authService inbound.AuthHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetString("username")
		if username == "" {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"error": "username not found in context"},
			)
			return
		}

		var req dto.ChangePassword
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(
				utils.GetBindJSONErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		req.Username = username
		err := authService.ChangePassword(c.Request.Context(), req.ToDomain())
		if err != nil {
			c.AbortWithStatusJSON(
				utils.GetDomainErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// Authenticate godoc
// @Summary Authenticate user
// @Description Authenticates the administrator and returns a token.
// @Tags Authentication
// @Accept x-www-form-urlencoded
// @Produce json
// @Param username formData string true "Login username"
// @Param password formData string true "Login password"
// @Success 200 {object} dto.AuthenticateResponse
// @Failure default {object} utils.ErrorResponse "Default error response for all failures"
// @Router /api/v1/auth/login [post]
func Authenticate(authService inbound.AuthHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.Authenticate

		if err := c.ShouldBind(&req); err != nil {
			c.AbortWithStatusJSON(
				utils.GetBindJSONErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		res, err := authService.Authenticate(c.Request.Context(), req.ToDomain())
		if err != nil {
			c.AbortWithStatusJSON(
				utils.GetDomainErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		c.JSON(http.StatusOK, dto.AuthenticateResponseFromDomain(res))
	}
}
