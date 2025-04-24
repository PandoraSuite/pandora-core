package entities

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type Service struct {
	ID int

	Name    string
	Status  enums.ServiceStatus
	Version string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *Service) Validate() *errors.Error {
	if s.Name == "" {
		return errors.ErrNameCannotBeEmpty
	}

	if s.Version == "" {
		return errors.ErrVersionCannotBeEmpty
	}

	return nil
}

func (a *Service) IsActive() bool {
	return a.Status == enums.ServiceActive
}
