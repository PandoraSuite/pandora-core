package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/MAD-py/pandora-core/internal/adapters/http/dto"
	"github.com/MAD-py/pandora-core/internal/adapters/http/errors"
	apikey "github.com/MAD-py/pandora-core/internal/app/api_key"
)

// APIKeyCreate godoc
// @Summary Creates a new API Key
// @Description Generates an API Key for a specific environment
// @Tags API Keys
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param request body dto.APIKeyCreate true "API Key creation data"
// @Success 201 {object} dto.APIKeyResponse
// @Failure default {object} errors.HTTPError "Default error response for all failures"
// @Router /api/v1/api-keys [post]
func APIKeyCreate(useCase apikey.CreateUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.APIKeyCreate

		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(errors.BindJSONToHTTPError(req, err))
			return
		}

		apiKey, err := useCase.Execute(c.Request.Context(), req.ToDomain())
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusCreated, dto.APIKeyResponseFromDomain(apiKey))
	}
}

// APIKeyUpdate godoc
// @Summary Updates an API key
// @Description Modifies the expiration date of a specific API key by ID
// @Tags API Keys
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param id path int true "API Key ID"
// @Param request body dto.APIKeyUpdate true "Updated API key data"
// @Success 200 {object} dto.APIKeyResponse
// @Failure default {object} errors.HTTPError "Default error response for all failures"
// @Router /api/v1/api-keys/{id} [patch]
func APIKeyUpdate(useCase apikey.UpdateUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKeyID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.Error(
				errors.NewValidationFailed(
					"path", "id", "Invalid api key id",
				),
			)
			return
		}

		var req dto.APIKeyUpdate
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(errors.BindJSONToHTTPError(req, err))
			return
		}

		apiKey, err := useCase.Execute(c.Request.Context(), apiKeyID, req.ToDomain())
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, dto.APIKeyResponseFromDomain(apiKey))
	}
}

// APIKeyDelete godoc
// @Summary Deletes an API key
// @Description Deletes a specific API key by ID
// @Tags API Keys
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param id path int true "API Key ID"
// @Success 204 "No Content"
// @Failure default {object} errors.HTTPError "Default error response for all failures"
// @Router /api/v1/api-keys/{id} [delete]
func APIKeyDelete(useCase apikey.DeleteUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKeyID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.Error(
				errors.NewValidationFailed(
					"path", "id", "Invalid api key id",
				),
			)
			return
		}

		if err := useCase.Execute(c.Request.Context(), apiKeyID); err != nil {
			c.Error(err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// APIKeyDisable godoc
// @Summary Disables an API key
// @Description Disables a specific API key by ID
// @Tags API Keys
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param id path int true "API Key ID"
// @Success 204 "No Content"
// @Failure default {object} errors.HTTPError "Default error response for all failures"
// @Router /api/v1/api-keys/{id}/disable [post]
func APIKeyDisable(useCase apikey.DisableUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKeyID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.Error(
				errors.NewValidationFailed(
					"path", "id", "Invalid api key id",
				),
			)
			return
		}

		if err := useCase.Execute(c.Request.Context(), apiKeyID); err != nil {
			c.Error(err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// APIKeyEnables godoc
// @Summary Enables an API key
// @Description Enables a specific API key by ID
// @Tags API Keys
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param id path int true "API Key ID"
// @Success 204 "No Content"
// @Failure default {object} errors.HTTPError "Default error response for all failures"
// @Router /api/v1/api-keys/{id}/enable [post]
func APIKeyEnable(useCase apikey.EnableUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKeyID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.Error(
				errors.NewValidationFailed(
					"path", "id", "Invalid api key id",
				),
			)
			return
		}

		if err := useCase.Execute(c.Request.Context(), apiKeyID); err != nil {
			c.Error(err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// APIKeyRevealKey godoc
// @Summary Reveals the API Key
// @Description Retrieves the actual API Key value for a specific API key by ID
// @Tags API Keys
// @Security ScopedToken
// @Accept json
// @Produce json
// @Param id path int true "API Key ID"
// @Success 200 {object} dto.APIKeyRevealKeyResponse
// @Failure default {object} errors.HTTPError "Default error response for all failures"
// @Router /api/v1/api-keys/{id}/reveal/key [get]
func APIKeyRevealKey(useCase apikey.RevealKeyUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKeyID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.Error(
				errors.NewValidationFailed(
					"path", "id", "Invalid api key id",
				),
			)
			return
		}

		apiKey, err := useCase.Execute(c.Request.Context(), apiKeyID)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, dto.APIKeyRevealKeyResponseFromDomain(apiKey))
	}
}
