package handlers

import (
	"net/http"
	"strconv"

	"github.com/MAD-py/pandora-core/internal/ports/inbound"
	"github.com/gin-gonic/gin"
)

// GetProjectsByClient godoc
// @Summary Retrieves all projects for a specific client
// @Description Fetches a list of projects associated with a given client
// @Tags Projects
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {array} dto.ProjectResponse
// @Failure 400 {object} map[string]string "Invalid client ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/clients/{id}/projects [get]
func GetProjectsByClient(projectService inbound.ProjectHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
			return
		}

		projects, err := projectService.GetProjectsByClient(c.Request.Context(), clientID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, projects)
	}
}
