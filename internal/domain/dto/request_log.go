package dto

import "time"

type RequestLogStatus string

const (
	RequestLogSuccess      RequestLogStatus = "success"
	RequestLogFailed       RequestLogStatus = "failed"
	RequestLogUnauthorized RequestLogStatus = "unauthorized"
	RequestLogServerError  RequestLogStatus = "server error"
)

type RequestLogCreate struct {
	APIKey        string `json:"api_key"`
	ServiceID     int    `json:"service_id"`
	EnvironmentID int    `json:"environment_id"`
}

type RequestLogResponse struct {
	ID              int              `json:"id"`
	APIKey          string           `json:"api_key"`
	ServiceID       int              `json:"service_id"`
	RequestTime     time.Time        `json:"request_time"`
	EnvironmentID   int              `json:"environment_id"`
	ExecutionStatus RequestLogStatus `json:"execution_status"`
	CreatedAt       time.Time        `json:"created_at"`
}
