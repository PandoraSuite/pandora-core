package models

import (
	"fmt"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

var clientType = []string{"organization", "developer"}

type Client struct {
	ID        pgtype.Int4
	Type      pgtype.Text
	Name      pgtype.Text
	Email     pgtype.Text
	CreatedAt pgtype.Timestamptz
}

func (c *Client) ValidateModel() error {
	return c.validateType()
}

func (c *Client) validateType() error {
	if t, _ := c.Type.Value(); t != nil {
		if slices.Contains(clientType, t.(string)) {
			return nil
		}
	}

	return fmt.Errorf(
		"invalid status: must be %s", strings.Join(clientType, ", "),
	)
}
