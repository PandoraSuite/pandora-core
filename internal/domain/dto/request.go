package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

// ... Requests ...

type RequestFilter struct {
	RequestTimeTo   time.Time                    `name:"request_time_to" validate:"omitempty,gtetime=RequestTimeTo"`
	RequestTimeFrom time.Time                    `name:"request_time_from" validate:"omitempty"`
	ExecutionStatus enums.RequestExecutionStatus `name:"execution_status" validate:"omitempty,enums=success forwarded client_error service_error unauthorized quota_exceeded"`
}

type RequestCreate struct {
	StartPoint      string                       `name:"start_point"`
	APIKey          string                       `name:"api_key"`
	APIKeyID        int                          `name:"api_key_id"`
	ProjectName     string                       `name:"project_name"`
	ProjectID       int                          `name:"project_id"`
	EnvironmentName string                       `name:"environment_name"`
	EnvironmentID   int                          `name:"environment_id"`
	ServiceName     string                       `name:"service_name"`
	ServiceVersion  string                       `name:"service_version"`
	ServiceID       int                          `name:"service_id"`
	StatusCode      int                          `name:"status_code"`
	ExecutionStatus enums.RequestExecutionStatus `name:"execution_status"`
	RequestTime     time.Time                    `name:"request_time"`
	Path            string                       `name:"path"`
	Method          string                       `name:"method"`
	IPAddress       string                       `name:"ip_address"`
	Body            string                       `name:"body"`
	BodyContentType enums.RequestBodyContentType `name:"body_content_type"`
	Headers         string                       `name:"headers"`
	QueryParams     string                       `name:"query_params"`
}

// ... Responses ...

type RequestResponse struct {
	ID              string                       `name:"id"`
	StartPoint      string                       `name:"start_point"`
	APIKey          string                       `name:"api_key"`
	APIKeyID        int                          `name:"api_key_id"`
	ProjectName     string                       `name:"project_name"`
	ProjectID       int                          `name:"project_id"`
	EnvironmentName string                       `name:"environment_name"`
	EnvironmentID   int                          `name:"environment_id"`
	ServiceName     string                       `name:"service_name"`
	ServiceVersion  string                       `name:"service_version"`
	ServiceID       int                          `name:"service_id"`
	StatusCode      int                          `name:"status_code"`
	ExecutionStatus enums.RequestExecutionStatus `name:"execution_status"`
	RequestTime     time.Time                    `name:"request_time"`
	Path            string                       `name:"path"`
	Method          string                       `name:"method"`
	IPAddress       string                       `name:"ip_address"`
	CreateAt        time.Time                    `name:"created_at"`
}

type RequestDetailsReponse struct {
	Body            string                       `name:"body"`
	BodyContentType enums.RequestBodyContentType `name:"body_content_type"`
	Headers         string                       `name:"headers"`
	QueryParams     string                       `name:"query_params"`
}

type RequestCompleteResponse struct {
	RequestResponse
	RequestDetailsReponse
}
