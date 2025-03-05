package handlers

import (
	"net/http"
	"strconv"

	"github.com/MAD-py/pandora-core/internal/ports/inbound"
	"github.com/gin-gonic/gin"
)

// GetAPIKeysByEnvironment godoc
// @Summary Retrieves all API Keys for an environment
// @Description Returns a list of API Keys associated with a specific environment
// @Tags Environments
// @Produce json
// @Param environment_id path int true "Environment ID"
// @Success 200 {array} dto.APIKeyResponse
// @Failure 400 {object} map[string]string "Invalid environment ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/environments/{environment_id}/api-keys [get]
func GetAPIKeysByEnvironment(apiKeyService inbound.APIKeyHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		environmentID, err := strconv.Atoi(c.Param("environment_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid environment_id"})
			return
		}

		apiKeys, err := apiKeyService.GetAPIKeysByEnvironment(
			c.Request.Context(), environmentID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, apiKeys)
	}
}
