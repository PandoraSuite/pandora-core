package entities

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type APIKey struct {
	ID int

	Key           string
	Status        enums.APIKeyStatus
	LastUsed      time.Time
	ExpiresAt     time.Time
	EnvironmentID int

	CreatedAt time.Time
}
