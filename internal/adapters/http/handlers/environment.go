package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/MAD-py/pandora-core/internal/adapters/http/handlers/utils"
	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/ports/inbound"
)

// CreateEnvironment godoc
// @Summary Creates a new environment
// @Description Adds a new environment to the system
// @Tags Environments
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param request body dto.EnvironmentCreate true "Environment creation data"
// @Success 201 {object} dto.EnvironmentResponse
// @Failure default {object} utils.ErrorResponse "Default error response for all failures"
// @Router /api/v1/environments [post]
func CreateEnvironment(environmentUseCase inbound.EnvironmentHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.EnvironmentCreate

		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(
				utils.GetBindJSONErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		environment, err := environmentUseCase.Create(c.Request.Context(), &req)
		if err != nil {
			c.AbortWithStatusJSON(
				utils.GetDomainErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		c.JSON(http.StatusCreated, environment)
	}
}

// GetEnvironment godoc
// @Summary Retrieves an environment by ID
// @Description Fetches the details of a specific environment using its ID
// @Tags Environments
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param id path int true "Environment ID"
// @Success 200 {object} dto.EnvironmentResponse
// @Failure default {object} utils.ErrorResponse "Default error response for all failures"
// @Router /api/v1/environments/{id} [get]
func GetEnvironment(environmentUseCase inbound.EnvironmentHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		environmentID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": "invalid environment ID"},
			)
			return
		}

		environment, err := environmentUseCase.GetByID(
			c.Request.Context(), environmentID,
		)
		if err != nil {
			c.AbortWithStatusJSON(
				utils.GetDomainErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		c.JSON(http.StatusOK, environment)
	}
}

// GetAPIKeysByEnvironment godoc
// @Summary Retrieves all API Keys for an environment
// @Description Returns a list of API Keys associated with a specific environment
// @Tags Environments
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param id path int true "Environment ID"
// @Success 200 {array} dto.APIKeyResponse
// @Failure default {object} utils.ErrorResponse "Default error response for all failures"
// @Router /api/v1/environments/{id}/api-keys [get]
func GetAPIKeysByEnvironment(apiKeyService inbound.APIKeyHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		environmentID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": "invalid environment ID"},
			)
			return
		}

		apiKeys, err := apiKeyService.GetAPIKeysByEnvironment(
			c.Request.Context(), environmentID,
		)
		if err != nil {
			c.AbortWithStatusJSON(
				utils.GetDomainErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		c.JSON(http.StatusOK, apiKeys)
	}
}

// AssignServiceToEnvironment godoc
// @Summary Assigns a service to an environment
// @Description Associates a service with a given environment
// @Tags Environments
// @Security OAuth2Password
// @Accept json
// @Produce json
// @Param id path int true "Environment ID"
// @Param request body dto.EnvironmentService true "Service data"
// @Success 204
// @Failure default {object} utils.ErrorResponse "Default error response for all failures"
// @Router /api/v1/environments/{id}/services [post]
func AssignServiceToEnvironment(environmentUseCase inbound.EnvironmentHTTPPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		environmentID, paramErr := strconv.Atoi(c.Param("id"))
		if paramErr != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"error": "Invalid environment ID"},
			)
			return
		}

		var req dto.EnvironmentService
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(
				utils.GetBindJSONErrorStatusCode(err),
				gin.H{"error": err.Error()},
			)
			return
		}

		err := environmentUseCase.AssignService(
			c.Request.Context(), environmentID, &req,
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
