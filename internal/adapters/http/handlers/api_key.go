package handlers

import (
	"net/http"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/ports/inbound"
	"github.com/gin-gonic/gin"
)

// CreateAPIKey godoc
// @Summary Creates a new API Key
// @Description Generates an API Key for a specific environment
// @Tags API Keys
// @Accept json
// @Produce json
// @Param request body dto.APIKeyCreate true "API Key creation data"
// @Success 201 {object} dto.APIKeyResponse
// @Failure 400 {object} map[string]string "Invalid input data"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/api-keys [post]
func CreateAPIKey(apiKeyService inbound.APIKeyHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.APIKeyCreate

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		apiKey, err := apiKeyService.Create(c.Request.Context(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, apiKey)
	}
}
