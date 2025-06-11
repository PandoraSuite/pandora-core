package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

// ... Requests ...

type RequestFilter struct {
	RequestTimeTo time.Time `form:"request_time_to" format:"date-time" extensions:"x-timezone=utc"`

	RequestTimeFrom time.Time `form:"request_time_from" format:"date-time" extensions:"x-timezone=utc"`

	ExecutionStatus string `form:"execution_status" enums:"success,forwarded,client_error,server_error,unauthorized,quota_exceeded"`
}

func (r *RequestFilter) ToDomain() *dto.RequestFilter {
	return &dto.RequestFilter{
		RequestTimeTo:   r.RequestTimeTo,
		RequestTimeFrom: r.RequestTimeFrom,
		ExecutionStatus: enums.RequestExecutionStatus(r.ExecutionStatus),
	}
}

// ... Responses ...

type RequestAPIKeyResponse struct {
	ID int `json:"id" minimum:"1"`

	Key string `json:"key" validate:"required" minLength:"11" maxLength:"11"`
}

type RequestServiceResponse struct {
	ID int `json:"id" minimum:"1"`

	Name string `json:"name" validate:"required"`

	Version string `json:"version" validate:"required" maxLength:"25"`
}

type RequestEnvironmentResponse struct {
	ID int `json:"id" validate:"required" minimum:"1"`

	Name string `json:"name" validate:"required"`
}

type RequestProjectResponse struct {
	ID int `json:"id" validate:"required" minimum:"1"`

	Name string `json:"name" validate:"required"`
}

type RequestResponse struct {
	ID string `json:"id" validate:"required" format:"uuid"`

	StartPoint string `json:"start_point"`

	APIKey *RequestAPIKeyResponse `json:"api_key" validate:"required"`

	Project *RequestProjectResponse `json:"project"`

	Environment *RequestEnvironmentResponse `json:"environment"`

	Service *RequestServiceResponse `json:"service" validate:"required"`

	StatusCode int `json:"status_code"`

	ExecutionStatus string `json:"execution_status" validate:"required" enums:"success,forwarded,client_error,server_error,unauthorized,quota_exceeded"`

	UnauthorizedReason string `json:"unauthorized_reason" enums:"API_KEY_INVALID,QUOTA_EXCEEDED,API_KEY_EXPIRED,API_KEY_DISABLED,SERVICE_MISMATCH,ENVIRONMENT_MISMATCH,ENVIRONMENT_DISABLED"`

	RequestTime time.Time `json:"request_time" validate:"required" format:"date-time" extensions:"x-timezone=utc"`

	Path string `json:"path" validate:"required"`

	Method string `json:"method" validate:"required" enums:"GET,HEAD,POST,PUT,PATCH,DELETE,CONNECT,OPTIONS,TRACE"`

	IPAddress string `json:"ip_address" validate:"required" format:"ipv4,ipv6"`

	CreateAt time.Time `json:"created_at" validate:"required" format:"date-time" extensions:"x-timezone=utc"`
}

func RequestResponseFromDomain(request *dto.RequestResponse) *RequestResponse {
	return &RequestResponse{
		ID:         request.ID,
		StartPoint: request.StartPoint,
		APIKey: &RequestAPIKeyResponse{
			ID:  request.APIKey.ID,
			Key: request.APIKey.Key,
		},
		Project: &RequestProjectResponse{
			ID:   request.Project.ID,
			Name: request.Project.Name,
		},
		Environment: &RequestEnvironmentResponse{
			ID:   request.Environment.ID,
			Name: request.Environment.Name,
		},
		Service: &RequestServiceResponse{
			ID:      request.Service.ID,
			Name:    request.Service.Name,
			Version: request.Service.Version,
		},
		StatusCode:         request.StatusCode,
		ExecutionStatus:    string(request.ExecutionStatus),
		UnauthorizedReason: string(request.UnauthorizedReason),
		RequestTime:        request.RequestTime,
		Path:               request.Path,
		Method:             request.Method,
		IPAddress:          request.IPAddress,
		CreateAt:           request.CreatedAt,
	}
}
