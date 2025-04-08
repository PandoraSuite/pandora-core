package entities

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/entities/utils"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type Client struct {
	ID int

	Type  enums.ClientType
	Name  string
	Email string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *Client) Validate() *errors.Error {
	if c.Type == enums.ClientTypeNull {
		return errors.ErrClientTypeCannotBeNull
	}

	if c.Name == "" {
		return errors.ErrNameCannotBeEmpty
	}

	if !utils.ValidateEmail(c.Email) {
		return errors.ErrInvalidEmailFormat
	}

	return nil
}
