package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type ProjectService struct {
	ID             int                                `json:"id" binding:"required"`
	MaxRequest     int                                `json:"max_request" binding:"required"`
	ResetFrequency enums.ProjectServiceResetFrequency `json:"reset_frequency" binding:"required" enums:",daily,weekly,biweekly,monthly" swaggertype:"string"`
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
	NextReset      time.Time                          `json:"next_reset"`
	MaxRequest     int                                `json:"max_request"`
	ResetFrequency enums.ProjectServiceResetFrequency `json:"reset_frequency"`
	AssignedAt     time.Time                          `json:"assigned_at"`
}

type ProjectResponse struct {
	ID        int                 `json:"id"`
	Name      string              `json:"name"`
	Status    enums.ProjectStatus `json:"status" enums:"in_production,in_development,deactivated" swaggertype:"string"`
	ClientID  int                 `json:"client_id"`
	CreatedAt time.Time           `json:"created_at"`

	Services []*ProjectServiceResponse `json:"services"`
}

type ProjectUpdate struct {
	Name string `json:"name,omitempty"`
}
