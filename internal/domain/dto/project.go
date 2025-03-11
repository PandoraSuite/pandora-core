package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type AssignServiceToProject struct {
	ProjectID      int                                `json:"-" swaggerignore:"true"`
	ServiceID      int                                `json:"-" swaggerignore:"true"`
	MaxRequest     int                                `json:"max_request"`
	ResetFrequency enums.ProjectServiceResetFrequency `json:"reset_frequency" enums:"null,daily,weekly,biweekly,monthly" swaggertype:"string"`
}

type ProjectCreate struct {
	Name     string              `json:"name"`
	Status   enums.ProjectStatus `json:"status" enums:"in_production,in_development" swaggertype:"string"`
	ClientID int                 `json:"client_id"`
}

type ProjectResponse struct {
	ID        int                 `json:"id"`
	Name      string              `json:"name"`
	Status    enums.ProjectStatus `json:"status" enums:"in_production,in_development,deactivated" swaggertype:"string"`
	ClientID  int                 `json:"client_id"`
	CreatedAt time.Time           `json:"created_at"`
}
