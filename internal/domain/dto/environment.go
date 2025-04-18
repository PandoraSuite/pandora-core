package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type DecrementAvailableRequest struct {
	MaxRequest       int `json:"max_request"`
	AvailableRequest int `json:"available_request"`
}

type QuotaUsage struct {
	MaxAllowed       int `json:"max_allowed"`
	CurrentAllocated int `json:"current_allocated"`
}

type EnvironmentService struct {
	ID         int `json:"id" binding:"required"`
	MaxRequest int `json:"max_request" binding:"required"`
}

type EnvironmentCreate struct {
	Name      string `json:"name"`
	ProjectID int    `json:"project_id"`

	Services []*EnvironmentService `json:"services,omitempty"`
}

type EnvironmentServiceResponse struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Version          string    `json:"version"`
	MaxRequest       int       `json:"max_request"`
	AvailableRequest int       `json:"available_request"`
	AssignedAt       time.Time `json:"assigned_at" time_format:"2006-01-02T15:04:05Z07:00" time_utc:"1"`
}

type EnvironmentResponse struct {
	ID        int                     `json:"id"`
	Name      string                  `json:"name"`
	Status    enums.EnvironmentStatus `json:"status" enums:"active,deactivated" swaggertype:"string"`
	ProjectID int                     `json:"project_id"`
	CreatedAt time.Time               `json:"created_at" time_format:"2006-01-02T15:04:05Z07:00" time_utc:"1"`

	Services []*EnvironmentServiceResponse `json:"services"`
}

type EnvironmentUpdate struct {
	Name string `json:"name,omitempty"`
}

type EnvironmentServiceUpdate struct {
	MaxRequest       int `json:"max_request"`
	AvailableRequest int `json:"-" swaggerignore:"true"`
}
