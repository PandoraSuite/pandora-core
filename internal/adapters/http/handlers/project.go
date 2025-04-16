package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/MAD-py/pandora-core/internal/adapters/http/handlers/utils"
	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/ports/inbound"
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

// GetProject godoc
// @Summary Retrieves a project by ID
// @Description Fetches the details of a specific project using its ID
// @Tags Projects
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} dto.ProjectResponse
// @Failure default {object} utils.ErrorResponse "Default error response for all failures"
// @Router /api/v1/projects/{id} [get]
func GetProject(projectService inbound.ProjectHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": "Invalid project ID"},
			)
			return
		}

		project, err := projectService.GetByID(
			c.Request.Context(), projectID,
		)
		if err != nil {
			c.AbortWithStatusJSON(
				utils.GetDomainErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		c.JSON(http.StatusOK, project)
	}
}

// UpdateProject godoc
// @Summary Updates a project
// @Description Modifies the details of a specific project by ID
// @Tags Projects
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param request body dto.ProjectUpdate true "Updated project data"
// @Success 200 {object} dto.ProjectResponse
// @Failure default {object} utils.ErrorResponse "Default error response for all failures"
// @Router /api/v1/projects/{id} [patch]
func UpdateProject(projectService inbound.ProjectHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": "Invalid Project ID"},
			)
			return
		}

		var req dto.ProjectUpdate
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(
				utils.GetBindJSONErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		project, err := projectService.Update(
			c.Request.Context(), projectID, &req,
		)
		if err != nil {
			c.AbortWithStatusJSON(
				utils.GetDomainErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		c.JSON(http.StatusOK, project)
	}
}

// GetEnvironmentsByProject godoc
// @Summary Retrieves all environments for a specific project
// @Description Fetches a list of environments associated with a given project
// @Tags Projects
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {array} dto.EnvironmentResponse
// @Failure default {object} utils.ErrorResponse "Default error response for all failures"
// @Router /api/v1/projects/{id}/environments [get]
func GetEnvironmentsByProject(projectService inbound.ProjectHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": "Invalid project ID"},
			)
			return
		}

		environments, err := projectService.GetEnvironments(
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

// AssignServiceToProject godoc
// @Summary Assigns a service to a project
// @Description Associates a service with a given project
// @Tags Projects
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param request body dto.ProjectService true "Service assignment data"
// @Success 204
// @Failure default {object} utils.ErrorResponse "Default error response for all failures"
// @Router /api/v1/projects/{id}/services [post]
func AssignServiceToProject(projectService inbound.ProjectHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": "Invalid project ID"},
			)
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

// RemoveServiceFromProject godoc
// @Summary Removes a service from a project
// @Description Disassociates a service from a specific project
// @Tags Projects
// @Security OAuth2Password
// @Produce json
// @Param id path int true "Project ID"
// @Param service_id path int true "Service ID"
// @Success 204
// @Failure default {object} utils.ErrorResponse "Default error response for all failures"
// @Router /api/v1/projects/{id}/services/{service_id} [delete]
func RemoveServiceFromProject(projectService inbound.ProjectHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": "Invalid project ID"},
			)
			return
		}

		serviceID, paramErr := strconv.Atoi(c.Param("service_id"))
		if paramErr != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": "Invalid service ID"},
			)
			return
		}

		err := projectService.RemoveService(
			c.Request.Context(), projectID, serviceID,
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
