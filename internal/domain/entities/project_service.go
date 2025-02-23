package entities

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type ProjectService struct {
	ProjectID int
	ServiceID int

	MaxRequest     int
	NextReset      time.Time
	ResetFrequency enums.ProjectServiceResetFrequency

	CreatedAt time.Time
}
