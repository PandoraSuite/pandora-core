package models

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence/models/utils"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
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

func (e *Environment) EntityID() int {
	return utils.PgtypeInt4ToInt(e.ID)
}

func (e *Environment) EntityCreatedAt() time.Time {
	return utils.PgtypeTimestamptzToTime(e.CreatedAt)
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

func (e *Environment) ToEntity() (*entities.Environment, error) {
	status, err := enums.ParseEnvironmentStatus(
		utils.PgtypeTextToString(e.Status),
	)
	if err != nil {
		return nil, err
	}

	return &entities.Environment{
		ID:        utils.PgtypeInt4ToInt(e.ID),
		Name:      utils.PgtypeTextToString(e.Name),
		Status:    status,
		ProjectID: utils.PgtypeInt4ToInt(e.ProjectID),
		CreatedAt: utils.PgtypeTimestamptzToTime(e.CreatedAt),
	}, nil
}

func EnvironmentToEntity(array []*Environment) ([]*entities.Environment, error) {
	result := make([]*entities.Environment, len(array))
	for i, v := range array {
		vv, err := v.ToEntity()
		if err != nil {
			return nil, err
		}
		result[i] = vv
	}
	return result, nil
}

func EnvironmentFromEntity(environment *entities.Environment) Environment {
	return Environment{
		ID:        utils.IntToPgtypeInt4(environment.ID),
		Name:      utils.StringToPgtypeText(environment.Name),
		Status:    utils.StringToPgtypeText(environment.Status.String()),
		ProjectID: utils.IntToPgtypeInt4(environment.ProjectID),
		CreatedAt: utils.TimeToPgtypeTimestamptz(environment.CreatedAt),
	}
}
