package entities

import (
	"time"
)

type Reservation struct {
	ID string

	EnvironmentID int
	ServiceID     int
	APIKey        string
	RequestTime   time.Time
	ExpiresAt     time.Time
}
