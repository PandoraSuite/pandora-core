package entities

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type Credentials struct {
	Username           string
	HashedPassword     string
	ForcePasswordReset bool
}

func (c *Credentials) CalculatePasswordHash(password string) errors.Error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return errors.NewInternal("unable to process the password", err)
	}

	c.HashedPassword = string(hashed)
	return nil
}

func (c *Credentials) VerifyPassword(password string) errors.Error {
	err := bcrypt.CompareHashAndPassword(
		[]byte(c.HashedPassword), []byte(password),
	)

	if err != nil {
		return errors.NewUnauthorized("invalid username or password")
	}
	return nil
}
