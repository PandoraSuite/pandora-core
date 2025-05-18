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
		status, paramErr := enums.ParseServiceStatus(c.Query("status"))
		if paramErr != nil {
			c.Error(
				errors.NewValidationFailed(
					"query", "type", "Invalid service status",
				),
			)
			return
		}

		req := dto.ServiceFilter{Status: status}
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
			c.Error(errors.BindingToDomainError(err))
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
			c.Error(errors.BindingToDomainError(err))
			return
		}

		service, err := useCase.Execute(
			c.Request.Context(), serviceID, req.Status,
		)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, dto.ServiceResponseFromDomain(service))
	}
}
