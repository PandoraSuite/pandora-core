package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

// ... Responses ...

type CheckStatusResponse struct {
	Status  enums.HealthStatus `json:"status" validate:"required"`
	Message string             `json:"message" validate:"required"`
	Latency int64              `json:"latency" validate:"required" minimum:"0"`
}

func CheckStatusResponseFromDomain(checkStatus *dto.CheckStatusResponse) *CheckStatusResponse {
	return &CheckStatusResponse{
		Status:  checkStatus.Status,
		Message: checkStatus.Message,
		Latency: checkStatus.Latency,
	}
}

type CheckResponse struct {
	Database *CheckStatusResponse `json:"database" validate:"required"`
}

func CheckResponseFromDomain(check *dto.CheckResponse) *CheckResponse {
	return &CheckResponse{
		Database: CheckStatusResponseFromDomain(check.Database),
	}
}

type HealthCheckResponse struct {
	Check     *CheckResponse     `json:"check" validate:"required"`
	Status    enums.HealthStatus `json:"status" validate:"required"`
	Timestamp time.Time          `json:"timestamp" validate:"required"`
}

func HealthCheckResponseFromDomain(healthCheck *dto.HealthCheckResponse) *HealthCheckResponse {
	return &HealthCheckResponse{
		Check:     CheckResponseFromDomain(healthCheck.Check),
		Status:    healthCheck.Status,
		Timestamp: healthCheck.Timestamp,
	}
}
