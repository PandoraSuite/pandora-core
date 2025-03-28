package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type APIKeyValidate struct {
	Key            string    `json:"key"`
	RequestTime    time.Time `json:"request_time"`
	ServiceName    string    `json:"service_name"`
	ServiceVersion string    `json:"service_version"`
}

type APIKeyValidateResponse struct {
	Valid     bool `json:"valid"`
	RequestID int  `json:"request_id"`

	Message string `json:"message,omitempty"`
}

type APIKeyValidateQuotaResponse struct {
	APIKeyValidateResponse `json:",inline"`

	AvailableRequest string `json:"available_request"`
}

type APIKeyCreate struct {
	ExpiresAt     time.Time `json:"expires_at"`
	EnvironmentID int       `json:"environment_id"`
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
