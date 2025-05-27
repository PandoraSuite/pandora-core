package entities

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type RequestMetadata struct {
	Body            string
	BodyContentType enums.RequestBodyContentType
	Headers         string
	QueryParams     string
}

type Request struct {
	ID string

	StartPoint      string
	APIKey          string
	APIKeyID        int
	ProjectName     string
	ProjectID       int
	EnvironmentName string
	EnvironmentID   int
	ServiceName     string
	ServiceVersion  string
	ServiceID       int
	StatusCode      int
	ExecutionStatus enums.RequestExecutionStatus
	RequestTime     time.Time
	Path            string
	Method          string
	IPAddress       string
	Metadata        *RequestMetadata

	CreatedAt time.Time
}
