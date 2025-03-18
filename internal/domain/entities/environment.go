package entities

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type Environment struct {
	ID int

	Name      string
	Status    enums.EnvironmentStatus
	ProjectID int

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
