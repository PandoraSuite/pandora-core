package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

// ... Requests ...

type APIKeyValidate struct {
	APIKey         string           `name:"api_key" validate:"required"`
	Request        *RequestIncoming `name:"request" validate:"required"`
	ServiceName    string           `name:"service_name" validate:"required"`
	ServiceVersion string           `name:"service_version" validate:"required"`
}

type APIKeyCreate struct {
	ExpiresAt     time.Time `name:"expires_at" validate:"omitempty,utc"`
	EnvironmentID int       `name:"environment_id" validate:"required,gt=0"`
}

type APIKeyUpdate struct {
	ExpiresAt time.Time `name:"expires_at" validate:"omitempty,utc"`
}

// ... Responses ...

type ConsumerInfo struct {
	ClientID        int    `name:"client_id"`
	ClientName      string `name:"client_name"`
	ProjectID       int    `name:"project_id"`
	ProjectName     string `name:"project_name"`
	EnvironmentID   int    `name:"environment_id"`
	EnvironmentName string `name:"environment_name"`
}

type APIKeyValidateResponse struct {
	Valid        bool                              `name:"valid"`
	RequestID    string                            `name:"request_id"`
	FailureCode  enums.APIKeyValidationFailureCode `name:"failure_code"`
	ConsumerInfo *ConsumerInfo                     `name:"consumer_info"`
}

type APIKeyValidateConsumeResponse struct {
	APIKeyValidateResponse
	AvailableRequest int `name:"available_request"`
}

type APIKeyResponse struct {
	ID            int                `name:"id"`
	Key           string             `name:"key"`
	Status        enums.APIKeyStatus `name:"status"`
	LastUsed      time.Time          `name:"last_used"`
	ExpiresAt     time.Time          `name:"expires_at"`
	EnvironmentID int                `name:"environment_id"`
	CreatedAt     time.Time          `name:"created_at"`
}
