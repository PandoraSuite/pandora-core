package entities

import (
	"fmt"
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
	if e.MaxRequest < -1 {
		return errors.ErrInvalidMaxRequest
	}

	if e.MaxRequest == -1 && e.AvailableRequest > -1 {
		return errors.ErrEnvironmentServiceAvailableRequestNotAllowed
	}

	if e.AvailableRequest > e.MaxRequest {
		return errors.ErrEnvironmentServiceAvailableRequestExceedsMax
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

	var errs []string
	for _, s := range e.Services {
		err := s.Validate()

		if err != nil {
			errs = append(
				errs,
				fmt.Sprintf("service %v: %s", s.ID, err.Message),
			)
		}
	}

	if len(errs) > 0 {
		return errors.NewError(
			errors.CodeValidationError,
			"invalid services assignments",
			errs...,
		)
	}

	return nil
}
