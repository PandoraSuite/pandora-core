package entities

import "time"

type ProjectService struct {
	ProjectID int
	ServiceID int

	MaxRequest     int
	NextReset      time.Time
	ResetFrequency string

	CreatedAt time.Time
}
