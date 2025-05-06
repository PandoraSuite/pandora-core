package dto

import "github.com/MAD-py/pandora-core/internal/domain/enums"

type ReservationWithDetails struct {
	ID                string                  `name:"id"`
	StartRequestID    string                  `name:"start_request_id"`
	APIKey            string                  `name:"api_key"`
	ServiceID         int                     `name:"service_id"`
	ServiceName       string                  `name:"service_name"`
	ServiceVersion    string                  `name:"service_version"`
	ServiceStatus     enums.ServiceStatus     `name:"service_status"`
	EnvironmentID     int                     `name:"environment_id"`
	EnvironmentName   string                  `name:"environment_name"`
	EnvironmentStatus enums.EnvironmentStatus `name:"environment_status"`
}
