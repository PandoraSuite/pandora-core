package entities

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/utils"
)

type Client struct {
	ID int

	Type  enums.ClientType
	Name  string
	Email string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *Client) Validate() errors.Error {
	if c.Type == enums.ClientTypeNull {
		return errors.NewValidationFailed(
			"client", "type", "type cannot be null",
		)
	}

	if c.Name == "" {
		return errors.NewValidationFailed(
			"client", "name", "name cannot be empty",
		)
	}

	if !utils.ValidateEmail(c.Email) {
		return errors.NewValidationFailed(
			"client", "email", "invalid email format",
		)
	}

	return nil
}
