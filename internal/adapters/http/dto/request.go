package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

// ... Requests ...

type RequestFilter struct {
	RequestTimeTo   time.Time `form:"request_time_to"`
	RequestTimeFrom time.Time `form:"request_time_from"`
	ExecutionStatus string    `form:"execution_status" enums:"success,forwarded,client_error,server_error,unauthorized,quota_exceeded"`
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
	ID                 string                      `json:"id"`
	StartPoint         string                      `json:"start_point"`
	APIKey             *RequestAPIKeyResponse      `json:"api_key"`
	Project            *RequestProjectResponse     `json:"project"`
	Environment        *RequestEnvironmentResponse `json:"environment"`
	Service            *RequestServiceResponse     `json:"service"`
	StatusCode         int                         `json:"status_code"`
	ExecutionStatus    string                      `json:"execution_status" enums:"success,forwarded,client_error,server_error,unauthorized,quota_exceeded"`
	UnauthorizedReason string                      `json:"unauthorized_reason" enums:"API_KEY_INVALID,QUOTA_EXCEEDED,API_KEY_EXPIRED,API_KEY_DISABLED,SERVICE_MISMATCH,ENVIRONMENT_MISMATCH,ENVIRONMENT_DISABLED"`
	RequestTime        time.Time                   `json:"request_time"`
	Path               string                      `json:"path"`
	Method             string                      `json:"method"`
	IPAddress          string                      `json:"ip_address"`
	CreateAt           time.Time                   `json:"created_at"`
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
