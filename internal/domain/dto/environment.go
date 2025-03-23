package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type DecrementAvailableRequest struct {
	MaxRequest       int `json:"max_request"`
	AvailableRequest int `json:"available_request"`
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
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Version    string    `json:"version"`
	MaxRequest int       `json:"max_request"`
	AssignedAt time.Time `json:"assigned_at"`
}

type EnvironmentResponse struct {
	ID        int                     `json:"id"`
	Name      string                  `json:"name"`
	Status    enums.EnvironmentStatus `json:"status" enums:"active,deactivated" swaggertype:"string"`
	ProjectID int                     `json:"project_id"`
	CreatedAt time.Time               `json:"created_at"`

	Services []*EnvironmentServiceResponse `json:"services"`
}

type EnvironmentUpdate struct {
	Name string `json:"name,omitempty"`
}
