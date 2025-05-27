package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
)

// ... Requests ...

type APIKeyCreate struct {
	ExpiresAt     time.Time `json:"expires_at"`
	EnvironmentID int       `json:"environment_id" binding:"required"`
}

func (a *APIKeyCreate) ToDomain() *dto.APIKeyCreate {
	return &dto.APIKeyCreate{
		ExpiresAt:     a.ExpiresAt,
		EnvironmentID: a.EnvironmentID,
	}
}

type APIKeyUpdate struct {
	ExpiresAt time.Time `json:"expires_at"`
}

func (a *APIKeyUpdate) ToDomain() *dto.APIKeyUpdate {
	return &dto.APIKeyUpdate{
		ExpiresAt: a.ExpiresAt,
	}
}

// ... Responses ...

type APIKeyResponse struct {
	ID            int       `json:"id"`
	Key           string    `json:"key"`
	Status        string    `json:"status" enums:"enabled,disabled,deprecated"`
	LastUsed      time.Time `json:"last_used" time_format:"2006-01-02T15:04:05Z07:00" time_utc:"1"`
	ExpiresAt     time.Time `json:"expires_at" time_format:"2006-01-02T15:04:05Z07:00" time_utc:"1"`
	EnvironmentID int       `json:"environment_id"`
	CreatedAt     time.Time `json:"created_at" time_format:"2006-01-02T15:04:05Z07:00" time_utc:"1"`
}

func APIKeyResponseFromDomain(apiKey *dto.APIKeyResponse) *APIKeyResponse {
	return &APIKeyResponse{
		ID:            apiKey.ID,
		Key:           apiKey.Key,
		Status:        string(apiKey.Status),
		LastUsed:      apiKey.LastUsed,
		ExpiresAt:     apiKey.ExpiresAt,
		EnvironmentID: apiKey.EnvironmentID,
		CreatedAt:     apiKey.CreatedAt,
	}
}
