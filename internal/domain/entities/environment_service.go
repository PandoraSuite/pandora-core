package entities

import "time"

type EnvironmentService struct {
	ServiceID     int
	EnvironmentID int

	MaxRequest       int
	AvailableRequest int

	CreatedAt time.Time
}
