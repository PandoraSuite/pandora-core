package entities

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
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

func (a *APIKey) GenerateKey() *errors.Error {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return errors.ErrAPIKeyGenerationFailed
	}

	a.Key = base64.RawURLEncoding.EncodeToString(bytes)
	return nil
}

func (a *APIKey) IsExpired() bool {
	return a.ExpiresAt.Before(time.Now())
}
