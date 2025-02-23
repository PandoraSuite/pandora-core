package dto

import "time"

type ProjectStatus string

const (
	ProjectInProduction  ProjectStatus = "in_production"
	ProjectInDevelopment ProjectStatus = "in_development"
	ProjectDeactivated   ProjectStatus = "deactivated"
)

type ProjectCreate struct {
	Name     string        `json:"name"`
	Status   ProjectStatus `json:"status"`
	ClientID int           `json:"client_id"`
}

type ProjectResponse struct {
	ID        int           `json:"id"`
	Name      string        `json:"name"`
	Status    ProjectStatus `json:"status"`
	ClientID  int           `json:"client_id"`
	CreatedAt time.Time     `json:"created_at"`
}
