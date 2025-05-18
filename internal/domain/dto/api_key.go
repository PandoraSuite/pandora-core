package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

// ... Requests ...

type APIKeyValidate struct {
	APIKey         string    `name:"api_key" validate:"required"`
	Service        string    `name:"service" validate:"required"`
	Environment    string    `name:"environment" validate:"required"`
	RequestTime    time.Time `name:"request_time" validate:"required,utc" time_format:"2006-01-02T15:04:05Z" time_utc:"1"`
	ServiceVersion string    `name:"service_version" validate:"required"`
}

type APIKeyValidateReserve struct {
	APIKeyValidate `validate:"required"`
	ReservationID  string `name:"reservation_id" validate:"required"`
}

type APIKeyCreate struct {
	ExpiresAt     time.Time `name:"expires_at" validate:"omitempty,utc"`
	EnvironmentID int       `name:"environment_id" validate:"required,gt=0"`
}

type APIKeyUpdate struct {
	ExpiresAt time.Time `name:"expires_at" validate:"omitempty,utc"`
}

// ... Responses ...

type APIKeyValidateResponse struct {
	Valid     bool   `name:"valid"`
	RequestID string `name:"request_id"`

	// Only when valid is true
	ReservationID    string `name:"reservation_id"` // Only for reservations
	AvailableRequest int    `name:"available_request"`

	// Only when valid is false
	Code    enums.ValidateStatusCode `name:"code"`
	Message string                   `name:"message"`
}

type APIKeyValidateConsumeResponse struct {
	APIKeyValidateResponse `name:",inline"`

	AvailableRequest string `name:"available_request"`
}

type APIKeyValidateReservationResponse struct {
	RequestID string                   `name:"request_id"`
	Valid     bool                     `name:"valid"`
	Message   string                   `name:"message"`
	Code      enums.ValidateStatusCode `name:"code"`
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

// ... Internal ...

type APIKeyEnabled struct {
	Code    enums.ValidateStatusCode
	Message string
}
