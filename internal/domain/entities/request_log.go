package entities

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type RequestLog struct {
	ID string

	APIKey          string
	ServiceID       int
	RequestTime     time.Time
	EnvironmentID   int
	StartPoint      string
	Message         string
	ExecutionStatus enums.RequestLogExecutionStatus

	CreatedAt time.Time
}
