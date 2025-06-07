package entities

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type EnvironmentService struct {
	ID int

	Name             string
	Version          string
	MaxRequests      int
	AvailableRequest int

	AssignedAt time.Time
}

func (es *EnvironmentService) Is(name, version string) bool {
	return es.Name == name && es.Version == version
}

type Environment struct {
	ID int

	Name      string
	Status    enums.EnvironmentStatus
	ProjectID int

	Services []*EnvironmentService

	CreatedAt time.Time
}

func (e *Environment) IsEnabled() bool {
	return e.Status == enums.EnvironmentStatusEnabled
}

func (e *Environment) Is(name string) bool {
	return e.Name == name
}
