package models

import (
	"fmt"
	"slices"
	"strings"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence/models/utils"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/jackc/pgx/v5/pgtype"
)

var projectStatus = []string{"in_production", "in_development", "deactivated"}

type Project struct {
	ID pgtype.Int4

	Name     pgtype.Text
	Status   pgtype.Text
	ClientID pgtype.Int4

	CreatedAt pgtype.Timestamptz
}

func (p *Project) ValidateModel() error {
	return p.validateStatus()
}

func (p *Project) validateStatus() error {
	if status, _ := p.Status.Value(); status != nil {
		if slices.Contains(projectStatus, status.(string)) {
			return nil
		}
	}

	return fmt.Errorf(
		"invalid status: must be %s", strings.Join(projectStatus, ", "),
	)
}

func (p *Project) ToEntity() *entities.Project {
	return &entities.Project{
		ID:        utils.PgtypeInt4ToInt(p.ID),
		Name:      utils.PgtypeTextToString(p.Name),
		Status:    utils.PgtypeTextToString(p.Status),
		ClientID:  utils.PgtypeInt4ToInt(p.ClientID),
		CreatedAt: utils.PgtypeTimestamptzToTime(p.CreatedAt),
	}
}

func ProjectFromEntity(project *entities.Project) *Project {
	return &Project{
		ID:        utils.IntToPgtypeInt4(project.ID),
		Name:      utils.StringToPgtypeText(project.Name),
		Status:    utils.StringToPgtypeText(project.Status),
		ClientID:  utils.IntToPgtypeInt4(project.ClientID),
		CreatedAt: utils.TimeToPgtypeTimestamptz(project.CreatedAt),
	}
}
