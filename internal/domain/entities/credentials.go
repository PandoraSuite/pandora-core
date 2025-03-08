package entities

import (
	"golang.org/x/crypto/bcrypt"

	domainErr "github.com/MAD-py/pandora-core/internal/domain/errors"
)

type Credentials struct {
	Username           string
	HashedPassword     string
	ForcePasswordReset bool
}

func (c *Credentials) CalculatePasswordHash(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return domainErr.ErrPasswordProcessingFailed
	}

	c.HashedPassword = string(hashed)
	return nil
}

func (c *Credentials) VerifyPassword(password string) error {
	err := bcrypt.CompareHashAndPassword(
		[]byte(c.HashedPassword), []byte(password),
	)

	if err != nil {
		return domainErr.ErrInvalidCredentials
	}
	return nil
}
