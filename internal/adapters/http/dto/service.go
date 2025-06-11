package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

// ... Requests ...

type ServiceFilter struct {
	Status string `form:"status" enums:"enabled,disabled,deprecated"`
}

func (s *ServiceFilter) ToDomain() *dto.ServiceFilter {
	return &dto.ServiceFilter{
		Status: enums.ServiceStatus(s.Status),
	}
}

type ServiceCreate struct {
	Name string `json:"name" validate:"required"`

	Version string `json:"version" validate:"required" maxLength:"25"`
}

func (s *ServiceCreate) ToDomain() *dto.ServiceCreate {
	return &dto.ServiceCreate{
		Name:    s.Name,
		Version: s.Version,
	}
}

type ServiceStatusUpdate struct {
	Status string `json:"status" validate:"required" enums:"enabled,disabled,deprecated"`
}

// ... Responses ...

type ServiceResponse struct {
	ID int `json:"id" validate:"required" minimum:"1"`

	Name string `json:"name" validate:"required"`

	Status string `json:"status" validate:"required" enums:"enabled,disabled,deprecated"`

	Version string `json:"version" validate:"required" maxLength:"25"`

	CreatedAt time.Time `json:"created_at" validate:"required" format:"date-time" extensions:"x-timezone=utc"`
}

func ServiceResponseFromDomain(service *dto.ServiceResponse) *ServiceResponse {
	return &ServiceResponse{
		ID:        service.ID,
		Name:      service.Name,
		Status:    string(service.Status),
		Version:   service.Version,
		CreatedAt: service.CreatedAt,
	}
}
