package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

// ... Requests ...

type EnvironmentService struct {
	ID         int `name:"id" validate:"required,gt=0"`
	MaxRequest int `name:"max_request" validate:"omitempty,gte=-1"`
}

type EnvironmentCreate struct {
	Name      string `name:"name" validate:"required"`
	ProjectID int    `name:"project_id" validate:"required,gt=0"`

	Services []*EnvironmentService `name:"services" validate:"required,dive"`
}

type EnvironmentUpdate struct {
	Name string `name:"name,omitempty"`
}

type EnvironmentServiceUpdateInput struct {
	MaxRequest int `name:"max_request" validate:"omitempty,gte=-1"`
}

// ... Responses ...

type EnvironmentServiceResponse struct {
	ID               int       `name:"id"`
	Name             string    `name:"name"`
	Version          string    `name:"version"`
	MaxRequest       int       `name:"max_request"`
	AvailableRequest int       `name:"available_request"`
	AssignedAt       time.Time `name:"assigned_at"`
}

type EnvironmentResponse struct {
	ID        int                     `name:"id"`
	Name      string                  `name:"name"`
	Status    enums.EnvironmentStatus `name:"status"`
	ProjectID int                     `name:"project_id"`
	CreatedAt time.Time               `name:"created_at"`

	Services []*EnvironmentServiceResponse `name:"services"`
}

type EnvironmentServiceReset struct {
	ID     int                     `name:"id"`
	Name   string                  `name:"name"`
	Status enums.EnvironmentStatus `name:"status"`

	Service *EnvironmentServiceResponse `name:"service"`
}

// ... Internal ...

type DecrementAvailableRequest struct {
	MaxRequest       int `name:"max_request"`
	AvailableRequest int `name:"available_request"`
}

type QuotaUsage struct {
	MaxAllowed       int `name:"max_allowed"`
	CurrentAllocated int `name:"current_allocated"`
}

type EnvironmentServiceUpdate struct {
	MaxRequest       int `name:"max_request"`
	AvailableRequest int `name:"available_request"`
}
