package entities

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type RequestMetadata struct {
	Body            string
	BodyContentType enums.RequestBodyContentType
	Cookies         string
	Headers         string
	QueryParams     string
}

type RequestAPIKey struct {
	ID  int
	Key string
}

func (r *RequestAPIKey) KeySummary() string {
	if len(r.Key) == 0 {
		return ""
	}

	return r.Key[:4] + "..." + r.Key[len(r.Key)-4:]
}

type RequestService struct {
	ID      int
	Name    string
	Version string
}

type RequestEnvironment struct {
	ID   int
	Name string
}

type RequestProject struct {
	ID   int
	Name string
}

type Request struct {
	ID string

	StartPoint         string
	APIKey             *RequestAPIKey
	Project            *RequestProject
	Environment        *RequestEnvironment
	Service            *RequestService
	Detail             string
	StatusCode         int
	ExecutionStatus    enums.RequestExecutionStatus
	UnauthorizedReason enums.APIKeyValidationFailureCode
	RequestTime        time.Time
	Path               string
	Method             string
	IPAddress          string
	Metadata           *RequestMetadata

	CreatedAt time.Time
}
