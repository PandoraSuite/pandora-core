package entities

import (
	"time"
)

type APIKey struct {
	ID int

	Key           string
	Status        string
	LastUsed      time.Time
	ExpiresAt     time.Time
	EnvironmentID int

	CreatedAt time.Time
}
