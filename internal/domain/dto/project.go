package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type ProjectService struct {
	ID             int                                `json:"id" binding:"required"`
	MaxRequest     int                                `json:"max_request" binding:"required"`
	ResetFrequency enums.ProjectServiceResetFrequency `json:"reset_frequency" enums:"daily,weekly,biweekly,monthly," swaggertype:"string"`
}

type ProjectCreate struct {
	Name     string              `json:"name" binding:"required"`
	Status   enums.ProjectStatus `json:"status" binding:"required" enums:"in_production,in_development" swaggertype:"string"`
	ClientID int                 `json:"client_id" binding:"required"`

	Services []*ProjectService `json:"services,omitempty"`
}

type ProjectServiceResponse struct {
	ID             int                                `json:"id"`
	Name           string                             `json:"name"`
	Version        string                             `json:"version"`
	NextReset      time.Time                          `json:"next_reset" time_format:"2006-01-02T15:04:05Z07:00" time_utc:"1"`
	MaxRequest     int                                `json:"max_request"`
	ResetFrequency enums.ProjectServiceResetFrequency `json:"reset_frequency" enums:"daily,weekly,biweekly,monthly," swaggertype:"string"`
	AssignedAt     time.Time                          `json:"assigned_at" time_format:"2006-01-02T15:04:05Z07:00" time_utc:"1"`
}

type ProjectResponse struct {
	ID        int                 `json:"id"`
	Name      string              `json:"name"`
	Status    enums.ProjectStatus `json:"status" enums:"in_production,in_development,deactivated" swaggertype:"string"`
	ClientID  int                 `json:"client_id"`
	CreatedAt time.Time           `json:"created_at" time_format:"2006-01-02T15:04:05Z07:00" time_utc:"1"`

	Services []*ProjectServiceResponse `json:"services"`
}

type ProjectUpdate struct {
	Name string `json:"name,omitempty"`
}

type ProjectServiceUpdate struct {
	NextReset      time.Time                          `json:"next_reset,omitempty" time_format:"2006-01-02T15:04:05Z07:00" time_utc:"1"`
	MaxRequest     int                                `json:"max_request"`
	ResetFrequency enums.ProjectServiceResetFrequency `json:"reset_frequency" enums:"daily,weekly,biweekly,monthly," swaggertype:"string"`
}

type ProjectServiceResetRequest struct {
	RecalculateNextReset bool `json:"recalculate_next_reset"`
}

type ProjectServiceResetRequestResponse struct {
	ResetCount          int                        `json:"reset_count"`
	ProjectService      *ProjectServiceResponse    `json:"project_service"`
	EnvironmentServices []*EnvironmentServiceReset `json:"environment_services"`
}
