package dto

import "time"

type EnvironmentCreate struct {
	Name      string `json:"name"`
	ProjectID int    `json:"project_id"`
}

type EnvironmentResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	ProjectID int       `json:"project_id"`
	CreatedAt time.Time `json:"created_at"`
}
