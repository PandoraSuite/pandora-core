package entities

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type Environment struct {
	ID int

	Name      string
	Status    enums.EnvironmentStatus
	ProjectID int

	CreatedAt time.Time
}
