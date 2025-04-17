package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/MAD-py/pandora-core/internal/adapters/http/handlers/utils"
	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/ports/inbound"
)

// GetAllServices godoc
// @Summary Retrieves all services
// @Description Fetches a list of all registered services
// @Tags Services
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param query query dto.ServiceFilter false "Query parameters"
// @Success 200 {array} dto.ServiceResponse
// @Failure default {object} utils.ErrorResponse "Default error response for all failures"
// @Router /api/v1/services [get]
func GetAllServices(srvService inbound.ServiceHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		status, paramErr := enums.ParseServiceStatus(c.Query("status"))
		if paramErr != nil {
			c.AbortWithStatusJSON(
				http.StatusUnprocessableEntity,
				utils.ErrorResponse{Error: paramErr},
			)
			return
		}

		req := dto.ServiceFilter{Status: status}
		services, err := srvService.GetServices(c.Request.Context(), &req)
		if err != nil {
			c.AbortWithStatusJSON(
				utils.GetBindJSONErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		c.JSON(http.StatusOK, services)
	}
}

// CreateService godoc
// @Summary Creates a new service
// @Description Adds a new service to the system
// @Tags Services
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param request body dto.ServiceCreate true "Service creation data"
// @Success 201 {object} dto.ServiceResponse
// @Failure default {object} utils.ErrorResponse "Default error response for all failures"
// @Router /api/v1/services [post]
func CreateService(srvService inbound.ServiceHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ServiceCreate

		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(
				utils.GetBindJSONErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		service, err := srvService.Create(c.Request.Context(), &req)
		if err != nil {
			c.AbortWithStatusJSON(
				utils.GetDomainErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		c.JSON(http.StatusCreated, service)
	}
}

// DeleteService godoc
// @Summary Deletes a service
// @Description Permanently removes a service by its ID
// @Tags Services
// @Security OAuth2Password
// @Produce json
// @Param id path int true "Service ID"
// @Success 204
// @Failure default {object} utils.ErrorResponse "Default error response for all failures"
// @Router /api/v1/services/{id} [delete]
func DeleteService(srvService inbound.ServiceHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": "Invalid Service ID"},
			)
			return
		}

		if err := srvService.Delete(c.Request.Context(), serviceID); err != nil {
			c.AbortWithStatusJSON(
				utils.GetDomainErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// UpdateStatus godoc
// @Summary Updates the status of a service
// @Description Changes the current status of a specific service by ID
// @Tags Services
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param id path int true "Service ID"
// @Param request body dto.ServiceStatusUpdate true "New service status"
// @Success 200 {object} dto.ServiceResponse
// @Failure default {object} utils.ErrorResponse "Default error response for all failures"
// @Router /api/v1/services/{id}/status [patch]
func UpdateStatusService(srvService inbound.ServiceHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": "Invalid Service ID"},
			)
			return
		}

		var req dto.ServiceStatusUpdate
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(
				utils.GetBindJSONErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		service, err := srvService.UpdateStatus(
			c.Request.Context(), serviceID, req.Status,
		)
		if err != nil {
			c.AbortWithStatusJSON(
				utils.GetDomainErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		c.JSON(http.StatusOK, service)
	}
}
