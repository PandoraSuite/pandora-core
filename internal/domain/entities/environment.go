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

func (e *EnvironmentService) Validate() errors.Error {
	if e.MaxRequest < -1 {
		return errors.NewValidationFailed(
			"environment.service",
			"max_request",
			"max request must be greater than or equal to -1",
		)
	}

	if e.MaxRequest == -1 && e.AvailableRequest > -1 {
		return errors.NewValidationFailed(
			"environment.service",
			"available_request",
			"available_request cannot be set when max_request is -1 (unlimited)",
		)
	}

	if e.AvailableRequest > e.MaxRequest {
		return errors.NewValidationFailed(
			"environment.service",
			"available_request",
			"available_request cannot be greater than max_request",
		)
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

func (e *Environment) Validate() errors.Error {
	if e.Name == "" {
		return errors.NewValidationFailed(
			"environment", "name", "name cannot be empty",
		)
	}

	errSet := errors.NewErrorSet()
	for _, s := range e.Services {
		if err := s.Validate(); err != nil {
			errSet.Add(err)
		}
	}

	if errSet.HasErrors() {
		return errSet
	}

	return nil
}

func (a *Environment) IsActive() bool {
	return a.Status == enums.EnvironmentActive
}
