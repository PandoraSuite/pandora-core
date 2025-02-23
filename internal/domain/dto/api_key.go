package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type APIKeyCreate struct {
	ExpiresAt     time.Time `json:"expires_at"`
	EnvironmentID int       `json:"environment_id"`
}

type APIKeyResponse struct {
	ID            int                `json:"id"`
	Key           string             `json:"key"`
	Status        enums.APIKeyStatus `json:"status"`
	LastUsed      time.Time          `json:"last_used"`
	ExpiresAt     time.Time          `json:"expires_at"`
	EnvironmentID int                `json:"environment_id"`
	CreatedAt     time.Time          `json:"created_at"`
}
