package models

import (
	"fmt"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

var projectStatus = []string{"in_production", "in_development", "deactivated"}

type Project struct {
	ID        pgtype.Int4
	ClientID  pgtype.Int4
	Name      pgtype.Text
	Status    pgtype.Text
	CreatedAt pgtype.Timestamptz
}

func (p *Project) ValidateModel() error {
	return p.validateStatus()
}

func (p *Project) validateStatus() error {
	if status, _ := p.Status.Value(); status != nil {
		if slices.Contains(projectStatus, status.(string)) {
			return nil
		}
	}

	return fmt.Errorf(
		"invalid status: must be %s", strings.Join(projectStatus, ", "),
	)
}
