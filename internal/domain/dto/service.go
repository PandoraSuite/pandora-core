package dto

import "time"

type ServiceStatus string

const (
	ServiceActive      ServiceStatus = "active"
	ServiceDeactivated ServiceStatus = "deactivated"
	ServiceDeprecated  ServiceStatus = "deprecated"
)

type ServiceCreate struct {
	Name    string        `json:"name"`
	Status  ServiceStatus `json:"status"`
	Version string        `json:"version"`
}

type ServiceResponse struct {
	ID        int           `json:"id"`
	Name      string        `json:"name"`
	Status    ServiceStatus `json:"status"`
	Version   string        `json:"version"`
	CreatedAt time.Time     `json:"created_at"`
}
