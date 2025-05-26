package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/MAD-py/pandora-core/internal/adapters/http/dto"
	"github.com/MAD-py/pandora-core/internal/adapters/http/errors"
	"github.com/MAD-py/pandora-core/internal/app/service"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

// ServiceList godoc
// @Summary Retrieves all services
// @Description Fetches a list of all registered services
// @Tags Services
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param query query dto.ServiceFilter false "Query parameters"
// @Success 200 {array} dto.ServiceResponse
// @Failure default {object} errors.HTTPError "Default error response for all failures"
// @Router /api/v1/services [get]
func ServiceList(useCase service.ListUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := dto.ServiceFilter{Status: c.Query("status")}
		services, err := useCase.Execute(c.Request.Context(), req.ToDomain())
		if err != nil {
			c.Error(err)
			return
		}

		resp := make([]*dto.ServiceResponse, len(services))
		for i, service := range services {
			resp[i] = dto.ServiceResponseFromDomain(service)
		}
		c.JSON(http.StatusOK, resp)
	}
}

// ServiceCreate godoc
// @Summary Creates a new service
// @Description Adds a new service to the system
// @Tags Services
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param request body dto.ServiceCreate true "Service creation data"
// @Success 201 {object} dto.ServiceResponse
// @Failure default {object} errors.HTTPError "Default error response for all failures"
// @Router /api/v1/services [post]
func ServiceCreate(useCase service.CreateUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ServiceCreate
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(errors.BindJSONToHTTPError(req, err))
			return
		}

		service, err := useCase.Execute(c.Request.Context(), req.ToDomain())
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusCreated, dto.ServiceResponseFromDomain(service))
	}
}

// ServiceDelete godoc
// @Summary Deletes a service
// @Description Permanently removes a service by its ID
// @Tags Services
// @Security OAuth2Password
// @Produce json
// @Param id path int true "Service ID"
// @Success 204
// @Failure default {object} errors.HTTPError "Default error response for all failures"
// @Router /api/v1/services/{id} [delete]
func ServiceDelete(useCase service.DeleteUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.Error(
				errors.NewValidationFailed(
					"path", "id", "Invalid service id",
				),
			)
			return
		}

		if err := useCase.Execute(c.Request.Context(), serviceID); err != nil {
			c.Error(err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// ServiceListRequests godoc
// @Summary Retrieves all requests for a service
// @Description Fetches a list of all requests associated with a specific service
// @Tags Services
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param id path int true "Service ID"
// @Param query query dto.RequestFilter false "Query parameters"
// @Success 200 {array} dto.RequestResponse
// @Failure default {object} errors.HTTPError "Default error response for all failures"
// @Router /api/v1/services/{id}/requests [get]
func ServiceListRequests(useCase service.ListRequestsUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.Error(
				errors.NewValidationFailed(
					"path", "id", "Invalid service id",
				),
			)
			return
		}

		var req dto.RequestFilter
		if err := c.ShouldBindQuery(&req); err != nil {
			c.Error(errors.BindQueryToHTTPError(req, err))
			return
		}

		requests, err := useCase.Execute(c.Request.Context(), serviceID, req.ToDomain())
		if err != nil {
			c.Error(err)
			return
		}

		resp := make([]*dto.RequestResponse, len(requests))
		for i, request := range requests {
			resp[i] = dto.RequestResponseFromDomain(request)
		}
		c.JSON(http.StatusOK, resp)
	}
}

// ServiceUpdateStatus godoc
// @Summary Updates the status of a service
// @Description Changes the current status of a specific service by ID
// @Tags Services
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param id path int true "Service ID"
// @Param request body dto.ServiceStatusUpdate true "New service status"
// @Success 200 {object} dto.ServiceResponse
// @Failure default {object} errors.HTTPError "Default error response for all failures"
// @Router /api/v1/services/{id}/status [patch]
func ServiceUpdateStatus(useCase service.UpdateStatusUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.Error(
				errors.NewValidationFailed(
					"path", "id", "Invalid service id",
				),
			)
			return
		}

		var req dto.ServiceStatusUpdate
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(errors.BindJSONToHTTPError(req, err))
			return
		}

		service, err := useCase.Execute(
			c.Request.Context(), serviceID, enums.ServiceStatus(req.Status),
		)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, dto.ServiceResponseFromDomain(service))
	}
}
