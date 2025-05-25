package entities

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
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

func (p *ProjectService) CalculateNextReset() {
	today := utils.TruncateToDay(time.Now())

	switch p.ResetFrequency {
	case enums.ProjectServiceResetFrequencyDaily:
		p.NextReset = today.AddDate(0, 0, 1)
	case enums.ProjectServiceResetFrequencyWeekly:
		p.NextReset = today.AddDate(0, 0, 7)
	case enums.ProjectServiceResetFrequencyBiweekly:
		p.NextReset = today.AddDate(0, 0, 14)
	case enums.ProjectServiceResetFrequencyMonthly:
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

func (p *Project) CalculateNextServicesReset() {
	for _, s := range p.Services {
		s.CalculateNextReset()
	}
}
