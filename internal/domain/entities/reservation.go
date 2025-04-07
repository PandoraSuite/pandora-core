package entities

import (
	"time"
)

type Reservation struct {
	ID string

	APIKey        string
	ServiceID     int
	EnvironmentID int
	RequestTime   time.Time
	ExpiresAt     time.Time
}
