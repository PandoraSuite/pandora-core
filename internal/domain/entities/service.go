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
