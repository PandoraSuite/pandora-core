package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
)

// ... Requests ...

type APIKeyCreate struct {
	// UTC
	ExpiresAt time.Time `json:"expires_at" format:"date-time"`

	EnvironmentID int `json:"environment_id" validate:"required" minimum:"1"`
}

func (a *APIKeyCreate) ToDomain() *dto.APIKeyCreate {
	return &dto.APIKeyCreate{
		ExpiresAt:     a.ExpiresAt,
		EnvironmentID: a.EnvironmentID,
	}
}

type APIKeyUpdate struct {
	// UTC
	ExpiresAt time.Time `json:"expires_at" format:"date-time"`
}

func (a *APIKeyUpdate) ToDomain() *dto.APIKeyUpdate {
	return &dto.APIKeyUpdate{
		ExpiresAt: a.ExpiresAt,
	}
}

// ... Responses ...

type APIKeyResponse struct {
	ID int `json:"id" minimum:"1"`

	Key string `json:"key" maxLength:"11" minLength:"11" example:"xxxx...xxxx"`

	Status string `json:"status" enums:"enabled,disabled,deprecated"`

	// UTC
	LastUsed time.Time `json:"last_used" format:"date-time"`

	// UTC
	ExpiresAt time.Time `json:"expires_at" format:"date-time"`

	EnvironmentID int `json:"environment_id" minimum:"1"`

	// UTC
	CreatedAt time.Time `json:"created_at" format:"date-time"`
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

type APIKeyRevealKeyResponse struct {
	Key string `json:"key" example:"xxxxxxxxxxxxx"`
}

func APIKeyRevealKeyResponseFromDomain(apiKey *dto.APIKeyRevealKeyResponse) *APIKeyRevealKeyResponse {
	return &APIKeyRevealKeyResponse{
		Key: apiKey.Key,
	}
}
