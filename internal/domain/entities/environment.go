package entities

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type EnvironmentService struct {
	ID int

	Name             string
	Version          string
	MaxRequest       int
	AvailableRequest int

	AssignedAt time.Time
}

func (e *EnvironmentService) Validate() *errors.Error {
	if e.ID <= 0 {
		return errors.ErrInvalidServiceID
	}

	if e.MaxRequest < 0 {
		return errors.ErrInvalidMaxRequest
	}

	return nil
}

type Environment struct {
	ID int

	Name      string
	Status    enums.EnvironmentStatus
	ProjectID int

	Services []*EnvironmentService

	CreatedAt time.Time
}

func (e *Environment) Validate() *errors.Error {
	if e.Name == "" {
		return errors.ErrNameCannotBeEmpty
	}

	if e.ProjectID <= 0 {
		return errors.ErrInvalidProjectID
	}

	return nil
}
