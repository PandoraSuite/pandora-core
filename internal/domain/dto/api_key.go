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

type APIKeyValidateReserve struct {
	APIKeyValidate `json:",inline"`
	ReservationID  string `json:"reservation_id"`
}
type APIKeyValidateBooking struct {
	Key       string `json:"key"`
	BookingID string `json:"booking_id"`
}

type APIKeyValidateResponse struct {
	RequestID        string                           `json:"request_id,omitempty"`
	AvailableRequest int                              `json:"available_request,omitempty"`
	Valid            bool                             `json:"valid"`
	Message          string                           `json:"message,omitempty"`
	Code             enums.ReserveExecutionStatusCode `json:"code,omitempty"`
}

type APIKeyValidateConsumeResponse struct {
	APIKeyValidateResponse `json:",inline"`

	AvailableRequest string `json:"available_request"`
}

type APIKeyValidateBookingResponse struct {
	APIKeyValidateConsumeResponse `json:",inline"`

	BookingID string `json:"booking_id"`
}

type APIKeyValidateReserveResponse struct {
	RequestID        string                           `json:"request_id,omitempty"`
	ReservationID    string                           `json:"reservation_id,omitempty"`
	AvailableRequest int                              `json:"available_request,omitempty"`
	Valid            bool                             `json:"valid"`
	Message          string                           `json:"message,omitempty"`
	Code             enums.ReserveExecutionStatusCode `json:"code,omitempty"`
}

type APIKeyValidateReservationResponse struct {
	RequestID string                           `json:"request_id,omitempty"`
	Valid     bool                             `json:"valid"`
	Message   string                           `json:"message,omitempty"`
	Code      enums.ReserveExecutionStatusCode `json:"code,omitempty"`
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
