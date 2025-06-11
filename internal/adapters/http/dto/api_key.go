package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
)

// ... Requests ...

type APIKeyCreate struct {
	ExpiresAt time.Time `json:"expires_at" format:"date-time" extensions:"x-timezone=utc"`

	EnvironmentID int `json:"environment_id" validate:"required" minimum:"1"`
}

func (a *APIKeyCreate) ToDomain() *dto.APIKeyCreate {
	return &dto.APIKeyCreate{
		ExpiresAt:     a.ExpiresAt,
		EnvironmentID: a.EnvironmentID,
	}
}

type APIKeyUpdate struct {
	ExpiresAt time.Time `json:"expires_at" format:"date-time" extensions:"x-timezone=utc"`
}

func (a *APIKeyUpdate) ToDomain() *dto.APIKeyUpdate {
	return &dto.APIKeyUpdate{
		ExpiresAt: a.ExpiresAt,
	}
}

// ... Responses ...

type APIKeyResponse struct {
	ID int `json:"id" validate:"required" minimum:"1"`

	Key string `json:"key" validate:"required" maxLength:"11" minLength:"11" example:"xxxx...xxxx"`

	Status string `json:"status" validate:"required" enums:"enabled,disabled,deprecated"`

	LastUsed time.Time `json:"last_used" format:"date-time" extensions:"x-timezone=utc"`

	ExpiresAt time.Time `json:"expires_at" format:"date-time" extensions:"x-timezone=utc"`

	EnvironmentID int `json:"environment_id" validate:"required" minimum:"1"`

	CreatedAt time.Time `json:"created_at" validate:"required" format:"date-time" extensions:"x-timezone=utc"`
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
	Key string `json:"key" validate:"required" example:"xxxxxxxxxxxxx"`
}

func APIKeyRevealKeyResponseFromDomain(apiKey *dto.APIKeyRevealKeyResponse) *APIKeyRevealKeyResponse {
	return &APIKeyRevealKeyResponse{
		Key: apiKey.Key,
	}
}
