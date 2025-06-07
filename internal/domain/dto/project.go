package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

// ... Requests ...

type ProjectService struct {
	ID             int                                `name:"id" validate:"required,gt=0"`
	MaxRequests    int                                `name:"max_requests" validate:"omitempty,gte=-1"`
	ResetFrequency enums.ProjectServiceResetFrequency `name:"reset_frequency" validate:"omitempty,enums=daily weekly biweekly monthly"`
}

type ProjectCreate struct {
	Name     string `name:"name" validate:"required"`
	ClientID int    `name:"client_id" validate:"required,gt=0"`

	Services []*ProjectService `name:"services" validate:"required,dive"`
}

type ProjectUpdate struct {
	Name string `name:"name" validate:"omitempty"`
}

type ProjectServiceUpdate struct {
	NextReset      time.Time                          `name:"next_reset" validate:"omitempty,utc"`
	MaxRequests    int                                `name:"max_requests" validate:"required,gte=-1"`
	ResetFrequency enums.ProjectServiceResetFrequency `name:"reset_frequency" validate:"omitempty,enums=daily weekly biweekly monthly"`
}

// ... Responses ...

type ProjectServiceResponse struct {
	ID             int                                `name:"id"`
	Name           string                             `name:"name"`
	Version        string                             `name:"version"`
	NextReset      time.Time                          `name:"next_reset"`
	MaxRequests    int                                `name:"max_requests"`
	ResetFrequency enums.ProjectServiceResetFrequency `name:"reset_frequency"`
	AssignedAt     time.Time                          `name:"assigned_at"`
}

type ProjectResponse struct {
	ID        int                 `name:"id"`
	Name      string              `name:"name"`
	Status    enums.ProjectStatus `name:"status"`
	ClientID  int                 `name:"client_id"`
	CreatedAt time.Time           `name:"created_at"`

	Services []*ProjectServiceResponse `name:"services"`
}

type ProjectResetRequestResponse struct {
	ResetCount          int                        `name:"reset_count"`
	ProjectService      *ProjectServiceResponse    `name:"project_service"`
	EnvironmentServices []*EnvironmentServiceReset `name:"environment_services"`
}

type ProjectClientInfoResponse struct {
	ProjectID   int    `name:"project_id"`
	ProjectName string `name:"project_name"`
	ClientID    int    `name:"client_id"`
	ClientName  string `name:"client_name"`
}
