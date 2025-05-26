package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

// ... Requests ...

type ServiceFilter struct {
	Status enums.ServiceStatus `form:"status" enums:"enabled,disabled,deprecated"`
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
	Status enums.ServiceStatus `json:"status" example:"enabled" enums:"enabled,disabled,deprecated"`
}

// ... Responses ...

type ServiceResponse struct {
	ID        int                 `json:"id"`
	Name      string              `json:"name"`
	Status    enums.ServiceStatus `json:"status" enums:"enabled,disabled,deprecated"`
	Version   string              `json:"version"`
	CreatedAt time.Time           `json:"created_at" time_format:"2006-01-02T15:04:05Z07:00" time_utc:"1"`
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
