package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

// ... Requests ...

type ServiceFilter struct {
	Status enums.ServiceStatus `form:"status,omitempty" enums:"active,deactivated,deprecated" swaggertype:"string"`
}

func (s *ServiceFilter) ToDomain() *dto.ServiceFilter {
	return &dto.ServiceFilter{
		Status: s.Status,
	}
}

type ServiceCreate struct {
	Name    string `json:"name" binding:"required"`
	Version string `json:"version" binding:"required"`
}

func (s *ServiceCreate) ToDomain() *dto.ServiceCreate {
	return &dto.ServiceCreate{
		Name:    s.Name,
		Version: s.Version,
	}
}

type ServiceStatusUpdate struct {
	Status enums.ServiceStatus `json:"status,omitempty" enums:"active,deactivated,deprecated" swaggertype:"string"`
}

// ... Responses ...

type ServiceResponse struct {
	ID        int                 `json:"id" example:"1"`
	Name      string              `json:"name" example:"Service Name"`
	Status    enums.ServiceStatus `json:"status" example:"active" enums:"active,deactivated,deprecated" swaggertype:"string"`
	Version   string              `json:"version" example:"1.0.0"`
	CreatedAt time.Time           `json:"created_at" example:"2025-01-01T00:00:00Z" time_format:"2006-01-02T15:04:05Z07:00" time_utc:"1"`
}

func ServiceResponseFromDomain(service *dto.ServiceResponse) *ServiceResponse {
	return &ServiceResponse{
		ID:        service.ID,
		Name:      service.Name,
		Status:    service.Status,
		Version:   service.Version,
		CreatedAt: service.CreatedAt,
	}
}
