package entities

import (
	"fmt"
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/utils"
)

type ProjectService struct {
	ID int

	Name           string
	Version        string
	NextReset      time.Time
	MaxRequest     int
	ResetFrequency enums.ProjectServiceResetFrequency

	AssignedAt time.Time
}

func (p *ProjectService) Validate() *errors.Error {
	if p.MaxRequest < -1 {
		return errors.ErrInvalidMaxRequest
	}

	if p.MaxRequest == -1 && p.ResetFrequency != enums.ProjectServiceNull {
		return errors.ErrProjectServiceResetFrequencyNotPermitted
	}

	if p.MaxRequest > -1 && p.ResetFrequency == enums.ProjectServiceNull {
		return errors.ErrProjectServiceResetFrequencyRequired
	}

	return nil
}

func (p *ProjectService) CalculateNextReset() {
	today := utils.TruncateToDay(time.Now())

	switch p.ResetFrequency {
	case enums.ProjectServiceDaily:
		p.NextReset = today.AddDate(0, 0, 1)
	case enums.ProjectServiceWeekly:
		p.NextReset = today.AddDate(0, 0, 7)
	case enums.ProjectServiceBiweekly:
		p.NextReset = today.AddDate(0, 0, 14)
	case enums.ProjectServiceMonthly:
		p.NextReset = today.AddDate(0, 1, 0)
	}
}

type Project struct {
	ID int

	Name     string
	Status   enums.ProjectStatus
	ClientID int

	Services []*ProjectService

	CreatedAt time.Time
}

func (p *Project) Validate() *errors.Error {
	if p.Status == enums.ProjectStatusNull {
		return errors.ErrProjectStatusCannotBeNull
	}

	if p.Name == "" {
		return errors.ErrNameCannotBeEmpty
	}

	var errs []string
	for _, s := range p.Services {
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

func (p *Project) CalculateNextServicesReset() {
	for _, s := range p.Services {
		s.CalculateNextReset()
	}
}

func (p *Project) GetService(serviceID int) (*ProjectService, *errors.Error) {
	for _, service := range p.Services {
		if service.ID == serviceID {
			return service, nil
		}
	}
	return nil, errors.ErrServiceNotFound
}
