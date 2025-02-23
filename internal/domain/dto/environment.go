package dto

import "time"

type EnvironmentStatus string

const (
	EnvironmentActive      EnvironmentStatus = "active"
	EnvironmentDeactivated EnvironmentStatus = "deactivated"
)

type EnvironmentCreate struct {
	Name      string `json:"name"`
	ProjectID int    `json:"project_id"`
}

type EnvironmentResponse struct {
	ID        int               `json:"id"`
	Name      string            `json:"name"`
	Status    EnvironmentStatus `json:"status"`
	ProjectID int               `json:"project_id"`
	CreatedAt time.Time         `json:"created_at"`
}
