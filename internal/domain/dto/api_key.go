package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type APIKeyValidate struct {
	Key            string    `json:"key"`
	Service        string    `json:"service"`
	Environment    string    `json:"environment"`
	RequestTime    time.Time `json:"request_time" time_format:"2006-01-02T15:04:05Z07:00" time_utc:"1"`
	ServiceVersion string    `json:"service_version"`
}

type APIKeyValidateReserve struct {
	APIKeyValidate `json:",inline"`
	ReservationID  string `json:"reservation_id"`
}
type APIKeyValidateBooking struct {
	Key       string `json:"key"`
	BookingID string `json:"booking_id"`
}

type APIKeyValidateResponse struct {
	RequestID        string                   `json:"request_id,omitempty"`
	ReservationID    string                   `json:"reservation_id,omitempty"` // Only for reservations
	AvailableRequest int                      `json:"available_request,omitempty"`
	ClientID         int                      `json:"client_id,omitempty"`
	Valid            bool                     `json:"valid"`
	Message          string                   `json:"message,omitempty"`
	Code             enums.ValidateStatusCode `json:"code,omitempty"`
}

type APIKeyValidateConsumeResponse struct {
	APIKeyValidateResponse `json:",inline"`

	AvailableRequest string `json:"available_request"`
}

type APIKeyValidateBookingResponse struct {
	APIKeyValidateConsumeResponse `json:",inline"`

	BookingID string `json:"booking_id"`
}

type APIKeyValidateReservationResponse struct {
	RequestID string                   `json:"request_id,omitempty"`
	Valid     bool                     `json:"valid"`
	Message   string                   `json:"message,omitempty"`
	Code      enums.ValidateStatusCode `json:"code,omitempty"`
}
type APIKeyCreate struct {
	ExpiresAt     time.Time `json:"expires_at" time_format:"2006-01-02T15:04:05Z07:00" time_utc:"1"`
	EnvironmentID int       `json:"environment_id" binding:"required"`
}

type APIKeyResponse struct {
	ID            int                `json:"id"`
	Key           string             `json:"key"`
	Status        enums.APIKeyStatus `json:"status" enums:"active,deactivated" swaggertype:"string"`
	LastUsed      time.Time          `json:"last_used" time_format:"2006-01-02T15:04:05Z07:00" time_utc:"1"`
	ExpiresAt     time.Time          `json:"expires_at" time_format:"2006-01-02T15:04:05Z07:00" time_utc:"1"`
	EnvironmentID int                `json:"environment_id"`
	CreatedAt     time.Time          `json:"created_at" time_format:"2006-01-02T15:04:05Z07:00" time_utc:"1"`
}

type APIKeyUpdate struct {
	ExpiresAt time.Time `json:"expires_at,omitempty" time_format:"2006-01-02T15:04:05Z07:00" time_utc:"1"`
}
