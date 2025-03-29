package handlers

import (
	"net/http"

	"github.com/MAD-py/pandora-core/internal/adapters/http/handlers/utils"
	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/ports/inbound"
	"github.com/gin-gonic/gin"
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
			c.JSON(
				utils.GetBindJSONErrorStatusCode(err),
				utils.ErrorResponse{Error: err},
			)
			return
		}

		apiKey, err := apiKeyService.Create(c.Request.Context(), &req)
		if err != nil {
			c.JSON(
				utils.GetDomainErrorStatusCode(err),
				utils.ErrorResponse{Error: err},
			)
			return
		}

		c.JSON(http.StatusCreated, apiKey)
	}
}
