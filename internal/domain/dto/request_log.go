package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type RequestLogCreate struct {
	APIKey          string                          `name:"api_key"`
	Message         string                          `name:"message"`
	ServiceID       int                             `name:"service_id"`
	RequestTime     time.Time                       `name:"request_time"`
	EnvironmentID   int                             `name:"environment_id"`
	ExecutionStatus enums.RequestLogExecutionStatus `name:"execution_status"`
}

type RequestLogResponse struct {
	ID              int                             `name:"id"`
	APIKey          string                          `name:"api_key"`
	ServiceID       int                             `name:"service_id"`
	RequestTime     time.Time                       `name:"request_time"`
	EnvironmentID   int                             `name:"environment_id"`
	ExecutionStatus enums.RequestLogExecutionStatus `name:"execution_status"`
	CreatedAt       time.Time                       `name:"created_at"`
}
