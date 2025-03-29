package handlers

import (
	"net/http"

	"github.com/MAD-py/pandora-core/internal/adapters/http/handlers/utils"
	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
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
// @Failure default {object} utils.ErrorResponse "Default error response for all failures"
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
// @Param query query dto.ServiceFilter false "Query parameters"
// @Success 200 {array} []dto.ServiceResponse
// @Failure default {object} utils.ErrorResponse "Default error response for all failures"
// @Router /api/v1/services [get]
func GetAllServices(srvService inbound.ServiceHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		s, paramErr := enums.ParseServiceStatus(c.Query("status"))
		if paramErr != nil {
			c.JSON(
				http.StatusUnprocessableEntity,
				utils.ErrorResponse{Error: paramErr},
			)
			return
		}

		req := dto.ServiceFilter{Status: s}
		services, err := srvService.GetServices(c.Request.Context(), &req)
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
