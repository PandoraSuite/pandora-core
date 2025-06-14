package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/MAD-py/pandora-core/internal/adapters/http/dto"
	"github.com/MAD-py/pandora-core/internal/app/health"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

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
