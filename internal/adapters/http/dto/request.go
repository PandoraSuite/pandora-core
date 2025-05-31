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
	ExecutionStatus string    `form:"execution_status" enums:"success,forwarded,client_error,service_error,unauthorized,quota_exceeded"`
}

func (r *RequestFilter) ToDomain() *dto.RequestFilter {
	return &dto.RequestFilter{
		RequestTimeTo:   r.RequestTimeTo,
		RequestTimeFrom: r.RequestTimeFrom,
		ExecutionStatus: enums.RequestExecutionStatus(r.ExecutionStatus),
	}
}

// ... Responses ...

type RequestResponse struct {
	ID              string    `json:"id"`
	StartPoint      string    `json:"start_point"`
	APIKey          string    `json:"api_key"`
	APIKeyID        int       `json:"api_key_id"`
	ProjectName     string    `json:"project_name"`
	ProjectID       int       `json:"project_id"`
	EnvironmentName string    `json:"environment_name"`
	EnvironmentID   int       `json:"environment_id"`
	ServiceName     string    `json:"service_name"`
	ServiceVersion  string    `json:"service_version"`
	ServiceID       int       `json:"service_id"`
	StatusCode      int       `json:"status_code"`
	ExecutionStatus string    `json:"execution_status" enums:"success,forwarded,client_error,service_error,unauthorized,quota_exceeded"`
	RequestTime     time.Time `json:"request_time"`
	Path            string    `json:"path"`
	Method          string    `json:"method"`
	IPAddress       string    `json:"ip_address"`
	CreateAt        time.Time `json:"created_at"`
}

func RequestResponseFromDomain(request *dto.RequestResponse) *RequestResponse {
	return &RequestResponse{
		ID:              request.ID,
		StartPoint:      request.StartPoint,
		APIKey:          request.APIKey,
		APIKeyID:        request.APIKeyID,
		ProjectName:     request.ProjectName,
		ProjectID:       request.ProjectID,
		EnvironmentName: request.EnvironmentName,
		EnvironmentID:   request.EnvironmentID,
		ServiceName:     request.ServiceName,
		ServiceVersion:  request.ServiceVersion,
		ServiceID:       request.ServiceID,
		StatusCode:      request.StatusCode,
		ExecutionStatus: string(request.ExecutionStatus),
		RequestTime:     request.RequestTime,
		Path:            request.Path,
		Method:          request.Method,
		IPAddress:       request.IPAddress,
		CreateAt:        request.CreateAt,
	}
}
