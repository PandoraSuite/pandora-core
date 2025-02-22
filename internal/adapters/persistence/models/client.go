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
	return c.validateStatus()
}

func (c *Client) validateStatus() error {
	if status, _ := c.Type.Value(); status != nil {
		if slices.Contains(clientType, status.(string)) {
			return nil
		}
	}

	return fmt.Errorf(
		"invalid status: must be %s", strings.Join(clientType, ", "),
	)
}
