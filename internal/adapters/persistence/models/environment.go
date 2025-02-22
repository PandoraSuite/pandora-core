package models

import (
	"fmt"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

var environmentStatus = []string{"active", "deactivated"}

type Environment struct {
	ID pgtype.Int4

	Name      pgtype.Text
	Status    pgtype.Text
	ProjectID pgtype.Int4

	CreatedAt pgtype.Timestamptz
}

func (e *Environment) ValidateModel() error {
	return e.validateStatus()
}

func (e *Environment) validateStatus() error {
	if status, _ := e.Status.Value(); status != nil {
		if slices.Contains(environmentStatus, status.(string)) {
			return nil
		}
	}

	return fmt.Errorf(
		"invalid status: must be %s", strings.Join(environmentStatus, ", "),
	)
}
