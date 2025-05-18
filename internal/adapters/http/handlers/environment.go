package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/MAD-py/pandora-core/internal/adapters/http/dto"
	"github.com/MAD-py/pandora-core/internal/adapters/http/errors"
	"github.com/MAD-py/pandora-core/internal/app/environment"
)

// EnvironmentCreate godoc
// @Summary Creates a new environment
// @Description Adds a new environment to the system
// @Tags Environments
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param request body dto.EnvironmentCreate true "Environment creation data"
// @Success 201 {object} dto.EnvironmentResponse
// @Failure default {object} errors.HTTPError "Default error response for all failures"
// @Router /api/v1/environments [post]
func EnvironmentCreate(useCase environment.CreateUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.EnvironmentCreate
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(errors.BindingToHTTPError(err))
			return
		}

		environment, err := useCase.Execute(c.Request.Context(), req.ToDomain())
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusCreated, dto.EnvironmentResponseFromDomain(environment))
	}
}

// EnvironmentGet godoc
// @Summary Retrieves an environment by ID
// @Description Fetches the details of a specific environment using its ID
// @Tags Environments
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param id path int true "Environment ID"
// @Success 200 {object} dto.EnvironmentResponse
// @Failure default {object} errors.HTTPError "Default error response for all failures"
// @Router /api/v1/environments/{id} [get]
func EnvironmentGet(useCase environment.GetUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		environmentID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.Error(
				errors.NewValidationFailed(
					"path", "id", "Invalid environment id",
				),
			)
			return
		}

		environment, err := useCase.Execute(c.Request.Context(), environmentID)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, dto.EnvironmentResponseFromDomain(environment))
	}
}

// EnvironmentUpdate godoc
// @Summary Updates an environment
// @Description Modifies the details of a specific environment by ID
// @Tags Environments
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param id path int true "Environment ID"
// @Param request body dto.EnvironmentUpdate true "Updated environment data"
// @Success 200 {object} dto.EnvironmentResponse
// @Failure default {object} errors.HTTPError "Default error response for all failures"
// @Router /api/v1/environments/{id} [patch]
func EnvironmentUpdate(useCase environment.UpdateUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		environmentID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.Error(
				errors.NewValidationFailed(
					"path", "id", "Invalid environment id",
				),
			)
			return
		}

		var req dto.EnvironmentUpdate
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(errors.BindingToHTTPError(err))
			return
		}

		environment, err := useCase.Execute(
			c.Request.Context(), environmentID, req.ToDomain(),
		)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, dto.EnvironmentResponseFromDomain(environment))
	}
}

// EnvironmentListAPIKeys godoc
// @Summary Retrieves all API Keys for an environment
// @Description Returns a list of API Keys associated with a specific environment
// @Tags Environments
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param id path int true "Environment ID"
// @Success 200 {array} dto.APIKeyResponse
// @Failure default {object} errors.HTTPError "Default error response for all failures"
// @Router /api/v1/environments/{id}/api-keys [get]
func EnvironmentListAPIKeys(useCase environment.ListAPIKeyUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		environmentID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.Error(
				errors.NewValidationFailed(
					"path", "id", "Invalid environment id",
				),
			)
			return
		}

		apiKeys, err := useCase.Execute(c.Request.Context(), environmentID)
		if err != nil {
			c.Error(err)
			return
		}

		resp := make([]*dto.APIKeyResponse, len(apiKeys))
		for i, apiKey := range apiKeys {
			resp[i] = dto.APIKeyResponseFromDomain(apiKey)
		}
		c.JSON(http.StatusOK, resp)
	}
}

// EnvironmentAssignService godoc
// @Summary Assigns a service to an environment
// @Description Associates a service with a given environment
// @Tags Environments
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param id path int true "Environment ID"
// @Param request body dto.EnvironmentService true "Service data"
// @Success 200 {object} dto.EnvironmentServiceResponse
// @Failure default {object} errors.HTTPError "Default error response for all failures"
// @Router /api/v1/environments/{id}/services [post]
func EnvironmentAssignService(useCase environment.AssignServiceUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		environmentID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.Error(
				errors.NewValidationFailed(
					"path", "id", "Invalid environment id",
				),
			)
			return
		}

		var req dto.EnvironmentService
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(errors.BindingToHTTPError(err))
			return
		}

		service, err := useCase.Execute(
			c.Request.Context(), environmentID, req.ToDomain(),
		)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, dto.EnvironmentServiceResponseFromDomain(service))
	}
}

// EnvironmentRemoveService godoc
// @Summary Removes a service from an environment
// @Description Disassociates a service from a specific environment
// @Tags Environments
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param id path int true "Environment ID"
// @Param service_id path int true "Service ID"
// @Success 204
// @Failure default {object} errors.HTTPError "Default error response for all failures"
// @Router /api/v1/environments/{id}/services/{service_id} [delete]
func EnvironmentRemoveService(useCase environment.RemoveServiceUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		environmentID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.Error(
				errors.NewValidationFailed(
					"path", "id", "Invalid environment id",
				),
			)
			return
		}

		serviceID, paramErr := strconv.Atoi(c.Param("service_id"))
		if paramErr != nil {
			c.Error(
				errors.NewValidationFailed(
					"path", "service_id", "Invalid service id",
				),
			)
			return
		}

		err := useCase.Execute(
			c.Request.Context(), environmentID, serviceID,
		)
		if err != nil {
			c.Error(err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// EnvironmentUpdateService godoc
// @Summary Updates a service assigned to an environment
// @Description Modifies the configuration of a service within a specific environment
// @Tags Environments
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param id path int true "Environment ID"
// @Param service_id path int true "Service ID"
// @Param request body dto.EnvironmentServiceUpdate true "Updated service configuration"
// @Success 200 {object} dto.EnvironmentServiceResponse
// @Failure default {object} errors.HTTPError "Default error response for all failures"
// @Router /api/v1/environments/{id}/services/{service_id} [patch]
func EnvironmentUpdateService(useCase environment.UpdateServiceUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		environmentID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.Error(
				errors.NewValidationFailed(
					"path", "id", "Invalid environment id",
				),
			)
			return
		}

		serviceID, paramErr := strconv.Atoi(c.Param("service_id"))
		if paramErr != nil {
			c.Error(
				errors.NewValidationFailed(
					"path", "service_id", "Invalid service id",
				),
			)
			return
		}

		var req dto.EnvironmentServiceUpdate
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(errors.BindingToHTTPError(err))
			return
		}

		service, err := useCase.Execute(
			c.Request.Context(), environmentID, serviceID, req.ToDomain(),
		)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, dto.EnvironmentServiceResponseFromDomain(service))
	}
}

// EnvironmentResetRequest godoc
// @Summary Resets request quota for a service in an environment
// @Description Resets the available request count for a specific service within an environment
// @Tags Environments
// @Security OAuth2Password
// @Produce json
// @Param id path int true "Environment ID"
// @Param service_id path int true "Service ID"
// @Success 200 {object} dto.EnvironmentServiceResponse
// @Failure default {object} errors.HTTPError "Default error response for all failures"
// @Router /api/v1/environments/{id}/services/{service_id}/reset-requests [patch]
func EnvironmentResetRequest(useCase environment.ResetRequestUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		environmentID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.Error(
				errors.NewValidationFailed(
					"path", "id", "Invalid environment id",
				),
			)
			return
		}

		serviceID, paramErr := strconv.Atoi(c.Param("service_id"))
		if paramErr != nil {
			c.Error(
				errors.NewValidationFailed(
					"path", "service_id", "Invalid service id",
				),
			)
			return
		}

		service, err := useCase.Execute(
			c.Request.Context(), environmentID, serviceID,
		)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, dto.EnvironmentServiceResponseFromDomain(service))
	}
}
