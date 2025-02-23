package entities

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type RequestLog struct {
	ID int

	APIKey          string
	ServiceID       int
	RequestTime     time.Time
	EnvironmentID   int
	ExecutionStatus enums.RequestLogExecutionStatus

	CreatedAt time.Time
}
