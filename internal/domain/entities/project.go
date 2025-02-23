package entities

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type Project struct {
	ID int

	Name     string
	Status   enums.ProjectStatus
	ClientID int

	CreatedAt time.Time
}
