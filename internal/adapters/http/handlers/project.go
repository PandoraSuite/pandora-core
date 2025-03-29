package handlers

import (
	"net/http"
	"strconv"

	"github.com/MAD-py/pandora-core/internal/adapters/http/handlers/utils"
	"github.com/MAD-py/pandora-core/internal/domain/dto"
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
// @Failure default {object} utils.ErrorResponse "Default error response for all failures"
// @Router /api/v1/projects [post]
func CreateProject(projectService inbound.ProjectHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ProjectCreate

		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(
				utils.GetBindJSONErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		project, err := projectService.Create(c.Request.Context(), &req)
		if err != nil {
			c.AbortWithStatusJSON(
				utils.GetDomainErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
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
// @Param id path int true "Project ID"
// @Param request body dto.ProjectService true "Service assignment data"
// @Success 204 "No Content"
// @Failure default {object} utils.ErrorResponse "Default error response for all failures"
// @Router /api/v1/projects/{id}/services [post]
func AssignServiceToProject(projectService inbound.ProjectHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
			return
		}

		var req dto.ProjectService
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": err.Error()},
			)
			return
		}

		err := projectService.AssignService(
			c.Request.Context(), projectID, &req,
		)
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

// GetEnvironmentsByProject godoc
// @Summary Retrieves all environments for a specific project
// @Description Fetches a list of environments associated with a given project
// @Tags Projects
// @Security OAuth2Password
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {array} dto.EnvironmentResponse
// @Failure default {object} utils.ErrorResponse "Default error response for all failures"
// @Router /api/v1/projects/{id}/environments [get]
func GetEnvironmentsByProject(environmentUseCase inbound.EnvironmentHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": "Invalid project ID"},
			)
			return
		}

		environments, err := environmentUseCase.GetByProject(
			c.Request.Context(), projectID,
		)
		if err != nil {
			c.AbortWithStatusJSON(
				utils.GetDomainErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		c.JSON(http.StatusOK, environments)
	}
}
