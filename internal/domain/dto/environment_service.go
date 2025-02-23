package dto

import "time"

type EnvironmentServiceCreate struct {
	ServiceID     int `json:"service_id"`
	EnvironmentID int `json:"environment_id"`
	MaxRequest    int `json:"max_request"`
}

type EnvironmentServiceResponse struct {
	ServiceID        int       `json:"service_id"`
	EnvironmentID    int       `json:"environment_id"`
	MaxRequest       int       `json:"max_request"`
	AvailableRequest int       `json:"available_request"`
	CreatedAt        time.Time `json:"created_at"`
}
