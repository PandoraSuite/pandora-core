package handlers

import (
	"net/http"
	"strconv"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	domainErr "github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/ports/inbound"
	"github.com/gin-gonic/gin"
)

// CreateProject godoc
// @Summary Creates a new project
// @Description Adds a new project to the system
// @Tags Projects
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param request body dto.ProjectCreate true "Project creation data"
// @Success 201 {object} dto.ProjectResponse
// @Failure 400 {object} map[string]string "Invalid input data"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /projects [post]
func CreateProject(projectService inbound.ProjectHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ProjectCreate

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		project, err := projectService.Create(c.Request.Context(), &req)
		if err != nil {
			if err == domainErr.ErrNameCannotBeEmpty {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		c.JSON(http.StatusCreated, project)
	}
}

// AssignServiceToProject godoc
// @Summary Assigns a service to a project
// @Description Associates a service with a given project
// @Tags Projects
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param project_id path int true "Project ID"
// @Param service_id path int true "Service ID"
// @Param request body dto.AssignServiceToProject true "Service assignment data"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "Invalid input data"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /projects/{project_id}/services/{service_id}/assign [post]
func AssignServiceToProject(projectService inbound.ProjectHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID, err := strconv.Atoi(c.Param("project_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
			return
		}

		serviceID, err := strconv.Atoi(c.Param("service_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
			return
		}

		var req dto.AssignServiceToProject
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		req.ProjectID = projectID
		req.ServiceID = serviceID

		err = projectService.AssignService(c.Request.Context(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// GetEnvironmentsByProject godoc
// @Summary Retrieves all environments for a specific project
// @Description Fetches a list of environments associated with a given project
// @Tags Environments
// @Security OAuth2Password
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {array} dto.EnvironmentResponse
// @Failure 400 {object} map[string]string "Invalid project ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /projects/{id}/environments [get]
func GetEnvironmentsByProject(environmentUseCase inbound.EnvironmentHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
			return
		}

		environments, err := environmentUseCase.GetEnvironmentsByProject(c.Request.Context(), projectID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, environments)
	}
}
