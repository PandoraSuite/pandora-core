package handlers

import (
	"net/http"

	"github.com/MAD-py/pandora-core/internal/adapters/http/handlers/utils"
	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/ports/inbound"
	"github.com/gin-gonic/gin"
)

// CreateService godoc
// @Summary Creates a new service
// @Description Adds a new service to the system
// @Tags Services
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param request body dto.ServiceCreate true "Service creation data"
// @Success 201 {object} dto.ServiceResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 422 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/services [post]
func CreateService(srvService inbound.ServiceHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ServiceCreate

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(
				utils.GetBindJSONErrorStatusCode(err),
				utils.ErrorResponse{Error: err},
			)
			return
		}

		service, err := srvService.Create(c.Request.Context(), &req)
		if err != nil {
			c.JSON(
				utils.GetDomainErrorStatusCode(err),
				utils.ErrorResponse{Error: err},
			)
			return
		}

		c.JSON(http.StatusCreated, service)
	}
}

// GetAllServices godoc
// @Summary Retrieves all services
// @Description Fetches a list of all registered services
// @Tags Services
// @Security OAuth2Password
// @Produce json
// @Success 200 {array} []dto.ServiceResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 422 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/services [get]
func GetAllServices(srvService inbound.ServiceHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		services, err := srvService.GetServices(c.Request.Context())
		if err != nil {
			c.JSON(
				utils.GetBindJSONErrorStatusCode(err),
				utils.ErrorResponse{Error: err},
			)
			return
		}

		c.JSON(http.StatusOK, services)
	}
}

// GetActiveServices godoc
// @Summary Retrieves active services
// @Description Fetches a list of all active services
// @Tags Services
// @Security OAuth2Password
// @Produce json
// @Success 200 {array} []dto.ServiceResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 422 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/services/active [get]
func GetActiveServices(srvService inbound.ServiceHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		services, err := srvService.GetActiveServices(c.Request.Context())
		if err != nil {
			c.JSON(
				utils.GetBindJSONErrorStatusCode(err),
				utils.ErrorResponse{Error: err},
			)
			return
		}

		c.JSON(http.StatusOK, services)
	}
}
