package dto

import "github.com/MAD-py/pandora-core/internal/domain/enums"

type ReservationWithDetails struct {
	ID                string                  `json:"id"`
	StartRequestID    string                  `json:"start_request_id"`
	APIKey            string                  `json:"api_key"`
	ServiceID         int                     `json:"service_id"`
	ServiceName       string                  `json:"service_name"`
	ServiceVersion    string                  `json:"service_version"`
	ServiceStatus     enums.ServiceStatus     `json:"service_status"`
	EnvironmentID     int                     `json:"environment_id"`
	EnvironmentName   string                  `json:"environment_name"`
	EnvironmentStatus enums.EnvironmentStatus `json:"environment_status"`
}
