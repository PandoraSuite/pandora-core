package dto

import "time"

type ProjectServiceStatus string

const (
	ProjectServiceDaily    ProjectServiceStatus = "daily"
	ProjectServiceWeekly   ProjectServiceStatus = "weekly"
	ProjectServiceBiweekly ProjectServiceStatus = "biweekly"
	ProjectServiceMonthly  ProjectServiceStatus = "monthly"
)

type ProjectServiceCreate struct {
	ProjectID      int                  `json:"project_id"`
	ServiceID      int                  `json:"service_id"`
	MaxRequest     string               `json:"max_request"`
	ResetFrequency ProjectServiceStatus `json:"reset_frequency"`
}

type ProjectServiceResponse struct {
	ProjectID      int                  `json:"project_id"`
	ServiceID      int                  `json:"service_id"`
	MaxRequest     string               `json:"max_request"`
	NextReset      time.Time            `json:"next_reset"`
	ResetFrequency ProjectServiceStatus `json:"reset_frequency"`
	CreatedAt      time.Time            `json:"created_at"`
}
