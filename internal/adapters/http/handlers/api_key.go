package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/MAD-py/pandora-core/internal/adapters/http/handlers/utils"
	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/ports/inbound"
)

// CreateAPIKey godoc
// @Summary Creates a new API Key
// @Description Generates an API Key for a specific environment
// @Tags API Keys
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param request body dto.APIKeyCreate true "API Key creation data"
// @Success 201 {object} dto.APIKeyResponse
// @Failure default {object} utils.ErrorResponse "Default error response for all failures"
// @Router /api/v1/api-keys [post]
func CreateAPIKey(apiKeyService inbound.APIKeyHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.APIKeyCreate

		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(
				utils.GetBindJSONErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		apiKey, err := apiKeyService.Create(c.Request.Context(), &req)
		if err != nil {
			c.AbortWithStatusJSON(
				utils.GetDomainErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		c.JSON(http.StatusCreated, apiKey)
	}
}

// UpdateAPIKey godoc
// @Summary Updates an API key
// @Description Modifies the expiration date of a specific API key by ID
// @Tags API Keys
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param id path int true "API Key ID"
// @Param request body dto.APIKeyUpdate true "Updated API key data"
// @Success 200 {object} dto.APIKeyResponse
// @Failure default {object} utils.ErrorResponse "Default error response for all failures"
// @Router /api/v1/api-keys/{id} [patch]
func UpdateAPIKey(apiKeyService inbound.APIKeyHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKeyID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": "Invalid API key ID"},
			)
			return
		}

		var req dto.APIKeyUpdate
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(
				utils.GetBindJSONErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		apiKey, err := apiKeyService.Update(c.Request.Context(), apiKeyID, &req)
		if err != nil {
			c.AbortWithStatusJSON(
				utils.GetDomainErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		c.JSON(http.StatusOK, apiKey)
	}
}
