package models

import (
	"fmt"
	"slices"
	"strings"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence/models/utils"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/jackc/pgx/v5/pgtype"
)

var environmentStatus = []string{"active", "deactivated"}

type Environment struct {
	ID pgtype.Int4

	Name      pgtype.Text
	Status    pgtype.Text
	ProjectID pgtype.Int4

	CreatedAt pgtype.Timestamptz
}

func (e *Environment) ValidateModel() error {
	return e.validateStatus()
}

func (e *Environment) validateStatus() error {
	if status, _ := e.Status.Value(); status != nil {
		if slices.Contains(environmentStatus, status.(string)) {
			return nil
		}
	}

	return fmt.Errorf(
		"invalid status: must be %s", strings.Join(environmentStatus, ", "),
	)
}

func (e *Environment) ToEntity() *entities.Environment {
	return &entities.Environment{
		ID:        utils.PgtypeInt4ToInt(e.ID),
		Name:      utils.PgtypeTextToString(e.Name),
		Status:    utils.PgtypeTextToString(e.Status),
		ProjectID: utils.PgtypeInt4ToInt(e.ProjectID),
		CreatedAt: utils.PgtypeTimestamptzToTime(e.CreatedAt),
	}
}

func EnvironmentFromEntity(environment *entities.Environment) *Environment {
	return &Environment{
		ID:        utils.IntToPgtypeInt4(environment.ID),
		Name:      utils.StringToPgtypeText(environment.Name),
		Status:    utils.StringToPgtypeText(environment.Status),
		ProjectID: utils.IntToPgtypeInt4(environment.ProjectID),
		CreatedAt: utils.TimeToPgtypeTimestamptz(environment.CreatedAt),
	}
}
