package entities

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type Service struct {
	ID int

	Name    string
	Status  enums.ServiceStatus
	Version string

	CreatedAt time.Time
}

func (s *Service) IsEnabled() bool {
	return s.Status == enums.ServiceStatusEnabled
}

func (s *Service) IsDisabled() bool {
	return s.Status == enums.ServiceStatusDisabled
}

func (s *Service) IsDeprecated() bool {
	return s.Status == enums.ServiceStatusDeprecated
}
