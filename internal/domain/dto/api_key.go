package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

// ... Requests ...

type APIKeyValidate struct {
	APIKey         string         `name:"api_key" validate:"required"`
	Request        *RequestCreate `name:"request" validate:"required"`
	ServiceName    string         `name:"service" validate:"required"`
	ServiceVersion string         `name:"service_version" validate:"required"`
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
	Valid        bool                              `name:"valid"`
	RequestID    string                            `name:"request_id"`
	FailureCode  enums.APIKeyValidationFailureCode `name:"failure_code"`
	ConsumerInfo *ProjectContextResponse           `name:"consumer_info"`
}

type APIKeyValidateConsumeResponse struct {
	APIKeyValidateResponse
	AvailableRequest int `name:"available_request"`
}

type APIKeyValidateResponse2 struct {
	Valid     bool   `name:"valid"`
	RequestID string `name:"request_id"`

	// Only when valid is true
	ReservationID    string `name:"reservation_id"` // Only for reservations
	AvailableRequest int    `name:"available_request"`

	// Only when valid is false
	Code    enums.ValidateStatusCode `name:"code"`
	Message string                   `name:"message"`
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
