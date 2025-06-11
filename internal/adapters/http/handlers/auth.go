package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/MAD-py/pandora-core/internal/adapters/http/dto"
	"github.com/MAD-py/pandora-core/internal/adapters/http/errors"
	"github.com/MAD-py/pandora-core/internal/app/auth"
)

// ChangePassword godoc
// @Summary Change password
// @Description Allows an authenticated user to change their password.
// @Tags Authentication
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param body body dto.ChangePassword true "Change password request"
// @Success 204
// @Failure default {object} errors.HTTPError "Default error response for all failures"
// @Router /api/v1/auth/change-password [post]
func ChangePassword(useCase auth.PasswordChangeUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetString("username")
		if username == "" {
			c.Error(errors.NewInternal("Username not found in context"))
			return
		}

		var req dto.ChangePassword
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(errors.BindJSONToHTTPError(req, err))
			return
		}

		req.Username = username
		err := useCase.Execute(c.Request.Context(), req.ToDomain())
		if err != nil {
			c.Error(err)
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
// @Param credentials formData dto.Authenticate true "Credenciales de acceso"
// @Success 200 {object} dto.AuthenticateResponse
// @Failure default {object} errors.HTTPError "Default error response for all failures"
// @Router /api/v1/auth/login [post]
func Authenticate(useCase auth.AutenticateUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.Authenticate
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errors.BindJSONToHTTPError(req, err))
			return
		}

		res, err := useCase.Execute(c.Request.Context(), req.ToDomain())
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, dto.AuthenticateResponseFromDomain(res))
	}
}

// Reauthenticate godoc
// @Summary Reauthenticate user
// @Description Reauthenticates the user for sensitive actions like revealing API keys.
// @Tags Authentication
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param body body dto.Reauthenticate true "Reauthentication request"
// @Success 200 {object} dto.ReauthenticateResponse
// @Failure default {object} errors.HTTPError "Default error response for all failures"
// @Router /api/v1/auth/reauthenticate [post]
func Reauthenticate(useCase auth.ReauthenticateUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetString("username")
		if username == "" {
			c.Error(errors.NewInternal("Username not found in context"))
			return
		}

		var req dto.Reauthenticate
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(errors.BindJSONToHTTPError(req, err))
			return
		}

		req.Username = username
		res, err := useCase.Execute(c.Request.Context(), req.ToDomain())
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, dto.ReauthenticateResponseFromDomain(res))
	}
}
