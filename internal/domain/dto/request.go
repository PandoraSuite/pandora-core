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

type RequestExecutionStatusUpdate struct {
	Detail          string                       `name:"detail" validate:"required_unless=execution_status success"`
	StatusCode      int                          `name:"status_code" validate:"required"`
	ExecutionStatus enums.RequestExecutionStatus `name:"execution_status" validate:"required,enums=success client_error service_error"`
}

// ... Responses ...

type RequestAPIKeyResponse struct {
	ID  int    `name:"id"`
	Key string `name:"key"`
}

type RequestServiceResponse struct {
	ID      int    `name:"id"`
	Name    string `name:"name"`
	Version string `name:"version"`
}

type RequestEnvironmentResponse struct {
	ID   int    `name:"id"`
	Name string `name:"name"`
}

type RequestProjectResponse struct {
	ID   int    `name:"id"`
	Name string `name:"name"`
}

type RequestResponse struct {
	ID                 string                            `name:"id"`
	StartPoint         string                            `name:"start_point"`
	APIKey             *RequestAPIKeyResponse            `name:"api_key"`
	Project            *RequestProjectResponse           `name:"project"`
	Environment        *RequestEnvironmentResponse       `name:"environment"`
	Service            *RequestServiceResponse           `name:"service"`
	Detail             string                            `name:"detail"`
	StatusCode         int                               `name:"status_code"`
	ExecutionStatus    enums.RequestExecutionStatus      `name:"execution_status"`
	UnauthorizedReason enums.APIKeyValidationFailureCode `name:"unauthorized_reason"`
	RequestTime        time.Time                         `name:"request_time"`
	Path               string                            `name:"path"`
	Method             string                            `name:"method"`
	IPAddress          string                            `name:"ip_address"`
	CreatedAt          time.Time                         `name:"created_at"`
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
