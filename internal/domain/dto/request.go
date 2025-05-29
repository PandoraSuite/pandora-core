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

type RequestIncomingMetadata struct {
	QueryParams     string                       `name:"query_params" validate:"omitempty"`
	Headers         string                       `name:"headers" validate:"omitempty"`
	Body            string                       `name:"body" validate:"omitempty"`
	BodyContentType enums.RequestBodyContentType `name:"body_content_type" validate:"omitempty,enums=application/xml application/json text/plain text/html multipart/form-data application/x-www-form-urlencoded application/octet-stream"`
}

type RequestIncoming struct {
	Path        string                   `name:"path" validate:"required"`
	Method      string                   `name:"method" validate:"required,enums=GET HEAD POST PUT PATCH DELETE CONNECT OPTIONS TRACE"`
	IPAddress   string                   `name:"ip_address" validate:"required,ip"`
	Metadata    *RequestIncomingMetadata `name:"metadata" validate:"omitempty"`
	RequestTime time.Time                `name:"request_time" validate:"required,utc" time_format:"2006-01-02T15:04:05Z" time_utc:"1"`
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
