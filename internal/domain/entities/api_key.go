package entities

import (
	"crypto/rand"
	"encoding/base64"
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

func (a *APIKey) GenerateKey() error {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return err
	}

	a.Key = base64.RawURLEncoding.EncodeToString(bytes)
	return nil
}

func (a *APIKey) IsExpired() bool {
	return a.ExpiresAt.Before(time.Now())
}
