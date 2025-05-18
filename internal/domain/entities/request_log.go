package entities

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type RequestLog struct {
	ID string

	APIKey          string
	Message         string
	ServiceID       int
	StartPoint      string
	RequestTime     time.Time
	EnvironmentID   int
	ExecutionStatus enums.RequestLogExecutionStatus

	CreatedAt time.Time
}
