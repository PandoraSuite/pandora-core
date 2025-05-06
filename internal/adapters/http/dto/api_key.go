package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

// ... Requests ...

type APIKeyCreate struct {
	ExpiresAt     time.Time `json:"expires_at,omitempty"`
	EnvironmentID int       `json:"environment_id" binding:"required"`
}

func (a *APIKeyCreate) ToDomain() *dto.APIKeyCreate {
	return &dto.APIKeyCreate{
		ExpiresAt:     a.ExpiresAt,
		EnvironmentID: a.EnvironmentID,
	}
}

type APIKeyUpdate struct {
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

func (a *APIKeyUpdate) ToDomain() *dto.APIKeyUpdate {
	return &dto.APIKeyUpdate{
		ExpiresAt: a.ExpiresAt,
	}
}

// ... Responses ...

type APIKeyResponse struct {
	ID            int                `json:"id" example:"1"`
	Key           string             `json:"key" example:"Wc4w-TpcKtfpLUE9Ve6U1hj1MK0y33qMIqXNgC2i4Ww"`
	Status        enums.APIKeyStatus `json:"status" example:"active" enums:"active,deactivated" swaggertype:"string"`
	LastUsed      time.Time          `json:"last_used" example:"2025-01-01T00:00:00Z" time_format:"2006-01-02T15:04:05Z07:00" time_utc:"1"`
	ExpiresAt     time.Time          `json:"expires_at" example:"2001-01-01T00:00:00Z" time_format:"2006-01-02T15:04:05Z07:00" time_utc:"1"`
	EnvironmentID int                `json:"environment_id" example:"1"`
	CreatedAt     time.Time          `json:"created_at" example:"2025-01-01T00:00:00Z" time_format:"2006-01-02T15:04:05Z07:00" time_utc:"1"`
}

func APIKeyResponseFromDomain(apiKey *dto.APIKeyResponse) *APIKeyResponse {
	return &APIKeyResponse{
		ID:            apiKey.ID,
		Key:           apiKey.Key,
		Status:        apiKey.Status,
		LastUsed:      apiKey.LastUsed,
		ExpiresAt:     apiKey.ExpiresAt,
		EnvironmentID: apiKey.EnvironmentID,
		CreatedAt:     apiKey.CreatedAt,
	}
}
