package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type RequestLogCreate struct {
	APIKey        string `json:"api_key"`
	ServiceID     int    `json:"service_id"`
	EnvironmentID int    `json:"environment_id"`
}

type RequestLogResponse struct {
	ID              int                             `json:"id"`
	APIKey          string                          `json:"api_key"`
	ServiceID       int                             `json:"service_id"`
	RequestTime     time.Time                       `json:"request_time" time_format:"2006-01-02T15:04:05Z07:00" time_utc:"1"`
	EnvironmentID   int                             `json:"environment_id"`
	ExecutionStatus enums.RequestLogExecutionStatus `json:"execution_status" enums:"success,failed,pending,unauthorized,server error" swaggertype:"string"`
	CreatedAt       time.Time                       `json:"created_at" time_format:"2006-01-02T15:04:05Z07:00" time_utc:"1"`
}
