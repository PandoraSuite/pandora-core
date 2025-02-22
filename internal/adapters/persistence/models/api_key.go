package models

import (
	"fmt"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

var apiKeyStatus = []string{"active", "deactivated"}

type APIKey struct {
	ID            pgtype.Int4
	EnvironmentID pgtype.Int4
	Key           pgtype.Text
	ExpiresAt     pgtype.Timestamptz
	LastUsed      pgtype.Timestamptz
	Status        pgtype.Text
	CreatedAt     pgtype.Timestamptz
}

func (s *APIKey) ValidateModel() error {
	return s.validateStatus()
}

func (s *APIKey) validateStatus() error {
	if status, _ := s.Status.Value(); status != nil {
		if slices.Contains(apiKeyStatus, status.(string)) {
			return nil
		}
	}

	return fmt.Errorf(
		"invalid status: must be %s", strings.Join(apiKeyStatus, ", "),
	)
}
