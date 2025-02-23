package dto

import "time"

type ProjectServiceCreate struct {
	ProjectID      int    `json:"project_id"`
	ServiceID      int    `json:"service_id"`
	MaxRequest     string `json:"max_request"`
	ResetFrequency string `json:"reset_frequency"`
}

type ProjectServiceResponse struct {
	ProjectID      int       `json:"project_id"`
	ServiceID      int       `json:"service_id"`
	MaxRequest     string    `json:"max_request"`
	NextReset      time.Time `json:"next_reset"`
	ResetFrequency string    `json:"reset_frequency"`
	CreatedAt      time.Time `json:"created_at"`
}
