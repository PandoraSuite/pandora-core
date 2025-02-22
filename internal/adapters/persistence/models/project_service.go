package models

import (
	"fmt"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

var projectServiceResetFrequency = []string{
	"daily", "weekly", "biweekly", "monthly",
}

type ProjectService struct {
	ProjectID      pgtype.Int4
	ServiceID      pgtype.Int4
	MaxRequest     pgtype.Int4
	ResetFrequency pgtype.Text
	NextReset      pgtype.Timestamptz
	CreatedAt      pgtype.Timestamptz
}

func (ps *ProjectService) ValidateModel() error {
	return ps.validateResetFrequency()
}

func (ps *ProjectService) validateResetFrequency() error {
	if resetFrequency, _ := ps.ResetFrequency.Value(); resetFrequency != nil {
		if slices.Contains(projectServiceResetFrequency, resetFrequency.(string)) {
			return nil
		}
	}

	return fmt.Errorf(
		"invalid reset frequency: must be %s",
		strings.Join(projectServiceResetFrequency, ", "),
	)
}
