package models

import (
	"fmt"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

var apiKeyStatus = []string{"active", "deactivated"}

type APIKey struct {
	ID pgtype.Int4

	Key           pgtype.Text
	Status        pgtype.Text
	LastUsed      pgtype.Timestamptz
	ExpiresAt     pgtype.Timestamptz
	EnvironmentID pgtype.Int4

	CreatedAt pgtype.Timestamptz
}

func (k *APIKey) ValidateModel() error {
	return k.validateStatus()
}

func (k *APIKey) validateStatus() error {
	if status, _ := k.Status.Value(); status != nil {
		if slices.Contains(apiKeyStatus, status.(string)) {
			return nil
		}
	}

	return fmt.Errorf(
		"invalid status: must be %s", strings.Join(apiKeyStatus, ", "),
	)
}
