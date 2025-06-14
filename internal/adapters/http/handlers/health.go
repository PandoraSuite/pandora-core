package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/MAD-py/pandora-core/internal/adapters/http/dto"
	"github.com/MAD-py/pandora-core/internal/app/health"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

// HealthCheck godoc
// @Summary Health Check
// @Description Check the health status of the application
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} dto.HealthCheckResponse
// @Failure 503 {object} dto.HealthCheckResponse "Service Unavailable"
// @Failure default {object} errors.HTTPError "Default error response for all failures"
// @Router /api/v1/health [get]
func HealthCheck(useCase health.CheckUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		healthCheck := useCase.Execute()

		var statusCode int
		if healthCheck.Status == enums.HealthStatusDown {
			statusCode = http.StatusServiceUnavailable
		} else {
			statusCode = http.StatusOK
		}

		c.JSON(
			statusCode,
			dto.HealthCheckResponseFromDomain(healthCheck),
		)
	}
}
