package models

import (
	"fmt"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

var serviceStatus = []string{"active", "deactivated", "deprecated"}

type Service struct {
	ID        pgtype.Int4
	Name      pgtype.Text
	Version   pgtype.Text
	Status    pgtype.Text
	CreatedAt pgtype.Timestamptz
}

func (s *Service) ValidateModel() error {
	return s.validateStatus()
}

func (s *Service) validateStatus() error {
	if status, _ := s.Status.Value(); status != nil {
		if slices.Contains(serviceStatus, status.(string)) {
			return nil
		}
	}

	return fmt.Errorf(
		"invalid status: must be %s",
		strings.Join(serviceStatus, ", "),
	)
}
