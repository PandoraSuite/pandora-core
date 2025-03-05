package handlers

import (
	"net/http"
	"strconv"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	domainErr "github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/ports/inbound"
	"github.com/gin-gonic/gin"
)

// GetAPIKeysByEnvironment godoc
// @Summary Retrieves all API Keys for an environment
// @Description Returns a list of API Keys associated with a specific environment
// @Tags Environments
// @Produce json
// @Param id path int true "Environment ID"
// @Success 200 {array} dto.APIKeyResponse
// @Failure 400 {object} map[string]string "Invalid environment ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/environments/{id}/api-keys [get]
func GetAPIKeysByEnvironment(apiKeyService inbound.APIKeyHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		environmentID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid environment ID"})
			return
		}

		apiKeys, err := apiKeyService.GetAPIKeysByEnvironment(
			c.Request.Context(), environmentID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, apiKeys)
	}
}

// CreateEnvironment godoc
// @Summary Creates a new environment
// @Description Adds a new environment to the system
// @Tags Environments
// @Accept json
// @Produce json
// @Param request body dto.EnvironmentCreate true "Environment creation data"
// @Success 201 {object} dto.EnvironmentResponse
// @Failure 400 {object} map[string]string "Invalid input data"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/environments [post]
func CreateEnvironment(environmentUseCase inbound.EnvironmentHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.EnvironmentCreate

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		environment, err := environmentUseCase.Create(c.Request.Context(), &req)
		if err != nil {
			if err == domainErr.ErrNameCannotBeEmpty {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		c.JSON(http.StatusCreated, environment)
	}
}

// AssignServiceToEnvironment godoc
// @Summary Assigns a service to an environment
// @Description Associates a service with a given environment
// @Tags Environments
// @Accept json
// @Produce json
// @Param environment_id path int true "Environment ID"
// @Param service_id path int true "Environment ID"
// @Param request body dto.AssignServiceToEnvironment true "Service assignment data"
// @Success 201 {object} dto.EnvironmentServiceResponse
// @Failure 400 {object} map[string]string "Invalid input data"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/environments/{environment_id}/services/{service_id}/assign [post]
func AssignServiceToEnvironment(environmentUseCase inbound.EnvironmentHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		environmentID, err := strconv.Atoi(c.Param("environment_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid environment ID"})
			return
		}

		serviceID, err := strconv.Atoi(c.Param("service_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
			return
		}

		var req dto.AssignServiceToEnvironment
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		req.ServiceID = serviceID
		req.EnvironmentID = environmentID

		err = environmentUseCase.AssignService(c.Request.Context(), &req)
		if err != nil {
			if err == domainErr.ErrMaxRequestExceededForServiceInProyect {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		c.Status(http.StatusNoContent)
	}
}
