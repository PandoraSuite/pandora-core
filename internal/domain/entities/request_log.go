package entities

import "time"

type RequestLog struct {
	ID int

	APIKey          string
	ServiceID       int
	RequestTime     time.Time
	EnvironmentID   int
	ExecutionStatus string

	CreatedAt time.Time
}
