package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type AssignServiceToProject struct {
	ProjectID      int                                `json:"-" swaggerignore:"true"`
	ServiceID      int                                `json:"-" swaggerignore:"true"`
	MaxRequest     int                                `json:"max_request"`
	ResetFrequency enums.ProjectServiceResetFrequency `json:"reset_frequency" enums:",daily,weekly,biweekly,monthly" swaggertype:"string"`
}

type ProjectServiceAssignment struct {
	ID             int                                `json:"id" binding:"required"`
	MaxRequest     int                                `json:"max_request" binding:"required"`
	ResetFrequency enums.ProjectServiceResetFrequency `json:"reset_frequency" binding:"required" enums:",daily,weekly,biweekly,monthly" swaggertype:"string"`
}

type ProjectCreate struct {
	Name     string              `json:"name" binding:"required"`
	Status   enums.ProjectStatus `json:"status" binding:"required" enums:"in_production,in_development" swaggertype:"string"`
	ClientID int                 `json:"client_id" binding:"required"`

	Services []*ProjectServiceAssignment `json:"services,omitempty"`
}

type ProjectServiceAssignmentResponse struct {
	ID             int                                `json:"id"`
	NextReset      time.Time                          `json:"next_reset"`
	MaxRequest     int                                `json:"max_request"`
	ResetFrequency enums.ProjectServiceResetFrequency `json:"reset_frequency"`
}

type ProjectResponse struct {
	ID        int                 `json:"id"`
	Name      string              `json:"name"`
	Status    enums.ProjectStatus `json:"status" enums:"in_production,in_development,deactivated" swaggertype:"string"`
	ClientID  int                 `json:"client_id"`
	CreatedAt time.Time           `json:"created_at"`

	Services []*ProjectServiceAssignmentResponse `json:"services"`
}
