package entities

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type Project struct {
	ID int

	Name     string
	Status   enums.ProjectStatus
	ClientID int

	CreatedAt time.Time
}

func (p *Project) Validate() error {
	if p.Status == enums.ProjectStatusNull {
		return errors.ErrProjectStatusCannotBeNull
	}

	if p.Name == "" {
		return errors.ErrNameCannotBeEmpty
	}

	return nil
}
