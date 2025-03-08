package entities

import (
	"golang.org/x/crypto/bcrypt"

	domainErr "github.com/MAD-py/pandora-core/internal/domain/errors"
)

type Credential struct {
	Username           string
	HashedPassword     string
	ForcePasswordReset bool
}

func (c *Credential) VerifyPassword(password string) error {
	err := bcrypt.CompareHashAndPassword(
		[]byte(c.HashedPassword), []byte(password),
	)

	if err != nil {
		return domainErr.ErrInvalidCredentials
	}
	return nil
}
