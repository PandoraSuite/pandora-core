package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type APIKeyValidate struct {
	Key            string    `json:"key"`
	Service        string    `json:"service"`
	Environment    string    `json:"environment"`
	RequestTime    time.Time `json:"request_time"`
	ServiceVersion string    `json:"service_version"`
}

type APIKeyValidateBooking struct {
	Key       string `json:"key"`
	BookingID string `json:"booking_id"`
}

type APIKeyValidateResponse struct {
	RequestID string `json:"request_id"`
}

type APIKeyValidateConsumeResponse struct {
	APIKeyValidateResponse `json:",inline"`

	AvailableRequest string `json:"available_request"`
}

type APIKeyValidateBookingResponse struct {
	APIKeyValidateConsumeResponse `json:",inline"`

	BookingID string `json:"booking_id"`
}

type APIKeyCreate struct {
	ExpiresAt     time.Time `json:"expires_at"`
	EnvironmentID int       `json:"environment_id" binding:"required"`
}

type APIKeyResponse struct {
	ID            int                `json:"id"`
	Key           string             `json:"key"`
	Status        enums.APIKeyStatus `json:"status" enums:"active,deactivated" swaggertype:"string"`
	LastUsed      time.Time          `json:"last_used"`
	ExpiresAt     time.Time          `json:"expires_at"`
	EnvironmentID int                `json:"environment_id"`
	CreatedAt     time.Time          `json:"created_at"`
}

type APIKeyUpdate struct {
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}
