package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/MAD-py/pandora-core/internal/adapters/http/dto"
	"github.com/MAD-py/pandora-core/internal/adapters/http/handlers/utils"
	"github.com/MAD-py/pandora-core/internal/app/project"
)

// ProjectList godoc
// @Summary Retrieves all projects
// @Description Fetches a complete list of projects in the system
// @Tags Projects
// @Security OAuth2Password
// @Produce json
// @Success 200 {array} dto.ProjectResponse
// @Failure default {object} utils.ErrorResponse "Default error response for all failures"
// @Router /api/v1/projects [get]
func ProjectList(useCase project.ListUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		projects, err := useCase.Execute(c.Request.Context())
		if err != nil {
			c.AbortWithStatusJSON(
				utils.GetDomainErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		resp := make([]*dto.ProjectResponse, len(projects))
		for i, project := range projects {
			resp[i] = dto.ProjectResponseFromDomain(project)
		}

		c.JSON(http.StatusOK, resp)
	}
}

// ProjectCreate godoc
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
func ProjectCreate(useCase project.CreateUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ProjectCreate

		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(
				utils.GetBindJSONErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		project, err := useCase.Execute(c.Request.Context(), req.ToDomain())
		if err != nil {
			c.AbortWithStatusJSON(
				utils.GetDomainErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		c.JSON(http.StatusCreated, dto.ProjectResponseFromDomain(project))
	}
}

// ProjectGet godoc
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
func ProjectGet(useCase project.GetUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": "Invalid project ID"},
			)
			return
		}

		project, err := useCase.Execute(c.Request.Context(), projectID)
		if err != nil {
			c.AbortWithStatusJSON(
				utils.GetDomainErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		c.JSON(http.StatusOK, dto.ProjectResponseFromDomain(project))
	}
}

// ProjectUpdate godoc
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
func ProjectUpdate(useCase project.UpdateUseCase) gin.HandlerFunc {
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

		project, err := useCase.Execute(
			c.Request.Context(), projectID, req.ToDomain(),
		)
		if err != nil {
			c.AbortWithStatusJSON(
				utils.GetDomainErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		c.JSON(http.StatusOK, dto.ProjectResponseFromDomain(project))
	}
}

// ProjectListEnvironments godoc
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
func ProjectListEnvironments(useCase project.ListEnvironmentsUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": "Invalid project ID"},
			)
			return
		}

		environments, err := useCase.Execute(c.Request.Context(), projectID)
		if err != nil {
			c.AbortWithStatusJSON(
				utils.GetDomainErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		resp := make([]*dto.EnvironmentResponse, len(environments))
		for i, environment := range environments {
			resp[i] = dto.EnvironmentResponseFromDomain(environment)
		}

		c.JSON(http.StatusOK, resp)
	}
}

// ProjectAssignService godoc
// @Summary Assigns a service to a project
// @Description Associates a service with a given project
// @Tags Projects
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param request body dto.ProjectService true "Service assignment data"
// @Success 200 {object} dto.ProjectServiceResponse
// @Failure default {object} utils.ErrorResponse "Default error response for all failures"
// @Router /api/v1/projects/{id}/services [post]
func ProjectAssignService(useCase project.AssignServiceUseCase) gin.HandlerFunc {
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

		service, err := useCase.Execute(
			c.Request.Context(), projectID, req.ToDomain(),
		)
		if err != nil {
			c.AbortWithStatusJSON(
				utils.GetDomainErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		c.JSON(http.StatusOK, dto.ProjectServiceResponseFromDomain(service))
	}
}

// ProjectRemoveService godoc
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
func ProjectRemoveService(useCase project.RemoveServiceUseCase) gin.HandlerFunc {
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

		err := useCase.Execute(c.Request.Context(), projectID, serviceID)
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

// ProjectUpdateService godoc
// @Summary Updates a service assigned to a project
// @Description Modifies the configuration of a service within a specific project
// @Tags Projects
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param service_id path int true "Service ID"
// @Param request body dto.ProjectServiceUpdate true "Updated service configuration"
// @Success 200 {object} dto.ProjectServiceResponse
// @Failure default {object} utils.ErrorResponse "Default error response for all failures"
// @Router /api/v1/projects/{id}/services/{service_id} [patch]
func ProjectUpdateService(useCase project.UpdateServiceUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": "Invalid Project ID"},
			)
			return
		}

		serviceID, paramErr := strconv.Atoi(c.Param("service_id"))
		if paramErr != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": "Invalid Service ID"},
			)
			return
		}

		var req dto.ProjectServiceUpdate
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(
				utils.GetBindJSONErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		service, err := useCase.Execute(
			c.Request.Context(), projectID, serviceID, req.ToDomain(),
		)
		if err != nil {
			c.AbortWithStatusJSON(
				utils.GetDomainErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		c.JSON(http.StatusOK, dto.ProjectServiceResponseFromDomain(service))
	}
}

// ProjectResetRequest godoc
// @Summary Resets available requests for a service in a project
// @Description Resets the request quota for a specific service assigned to a project
// @Tags Projects
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param service_id path int true "Service ID"
// @Param request body dto.ProjectServiceResetRequest true "Reset configuration"
// @Success 200 {object} dto.ProjectServiceResetRequestResponse
// @Failure default {object} utils.ErrorResponse "Default error response for all failures"
// @Router /api/v1/projects/{id}/services/{service_id}/reset-requests [post]
func ProjectResetRequest(useCase project.ResetRequestUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": "Invalid Project ID"},
			)
			return
		}

		serviceID, paramErr := strconv.Atoi(c.Param("service_id"))
		if paramErr != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": "Invalid Service ID"},
			)
			return
		}

		var req dto.ProjectResetRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(
				utils.GetBindJSONErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		resp, err := useCase.Execute(
			c.Request.Context(), projectID, serviceID, req.RecalculateNextReset,
		)
		if err != nil {
			c.AbortWithStatusJSON(
				utils.GetDomainErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		c.JSON(http.StatusOK, dto.ProjectResetRequestResponseFromDomain(resp))
	}
}
