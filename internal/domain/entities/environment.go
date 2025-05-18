package entities

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type EnvironmentService struct {
	ID int

	Name             string
	Version          string
	MaxRequest       int
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

func (e *Environment) IsActive() bool {
	return e.Status == enums.EnvironmentActive
}

func (e *Environment) Is(name string) bool {
	return e.Name == name
}
