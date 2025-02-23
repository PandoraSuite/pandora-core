package dto

import (
	"time"
)

type APIKeyCreate struct {
	ExpiresAt     time.Time `json:"expires_at"`
	EnvironmentID int       `json:"environment_id"`
}

type APIKeyResponse struct {
	ID            int       `json:"id"`
	Key           string    `json:"key"`
	Status        string    `json:"status"`
	LastUsed      time.Time `json:"last_used"`
	ExpiresAt     time.Time `json:"expires_at"`
	EnvironmentID int       `json:"environment_id"`
	CreatedAt     time.Time `json:"created_at"`
}
