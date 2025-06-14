package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

// ... Responses ...

type CheckStatusResponse struct {
	Status  enums.HealthStatus `name:"status"`
	Message string             `name:"message"`
	Latency int64              `name:"latency"`
}

type CheckResponse struct {
	Database *CheckStatusResponse `name:"database"`
}

type HealthCheckResponse struct {
	Check     *CheckResponse     `name:"check"`
	Status    enums.HealthStatus `name:"status"`
	Timestamp time.Time          `name:"timestamp"`
}
