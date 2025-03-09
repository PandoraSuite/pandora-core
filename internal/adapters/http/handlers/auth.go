package handlers

import (
	"net/http"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	domainErr "github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/ports/inbound"
	"github.com/gin-gonic/gin"
)

// Authenticate godoc
// @Summary Authenticate user
// @Description Authenticates the administrator and returns a token.
// @Tags Authentication
// @Accept x-www-form-urlencoded
// @Produce json
// @Param username formData string true "Login username"
// @Param password formData string true "Login password"
// @Success 200 {object} dto.AuthenticateResponse
// @Failure 400 {object} map[string]string "Invalid request format"
// @Failure 401 {object} map[string]string "Invalid username or password"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/auth/login [post]
func Authenticate(authService inbound.AuthHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.Authenticate

		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := authService.Authenticate(c.Request.Context(), &req)
		if err != nil {
			if err == domainErr.ErrInvalidCredentials {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

// ChangePassword godoc
// @Summary Change password
// @Description Allows an authenticated user to change their password.
// @Tags Authentication
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param request body dto.ChangePassword true "New password and confirmation"
// @Success 204
// @Failure 400 {object} map[string]string "Validation error"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/auth/change-password [post]
func ChangePassword(authService inbound.AuthHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetString("username")
		if username == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "username not found in context"})
			return
		}

		var req dto.ChangePassword
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		req.Username = username
		err := authService.ChangePassword(c.Request.Context(), &req)
		if err != nil {
			switch err {
			case domainErr.ErrPasswordTooShort, domainErr.ErrPasswordMismatch:
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		c.Status(http.StatusNoContent)
	}
}
