package models

import (
	"fmt"
	"slices"
	"strings"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence/models/utils"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/jackc/pgx/v5/pgtype"
)

var projectServiceResetFrequency = []string{
	"daily", "weekly", "biweekly", "monthly",
}

type ProjectService struct {
	ProjectID pgtype.Int4
	ServiceID pgtype.Int4

	MaxRequest     pgtype.Int4
	NextReset      pgtype.Timestamptz
	ResetFrequency pgtype.Text

	CreatedAt pgtype.Timestamptz
}

func (ps *ProjectService) ValidateModel() error {
	return ps.validateResetFrequency()
}

func (ps *ProjectService) validateResetFrequency() error {
	if resetFrequency, _ := ps.ResetFrequency.Value(); resetFrequency != nil {
		if slices.Contains(projectServiceResetFrequency, resetFrequency.(string)) {
			return nil
		}
	}

	return fmt.Errorf(
		"invalid reset frequency: must be %s",
		strings.Join(projectServiceResetFrequency, ", "),
	)
}

func (ps *ProjectService) ToEntity() *entities.ProjectService {
	return &entities.ProjectService{
		ProjectID:      utils.PgtypeInt4ToInt(ps.ProjectID),
		ServiceID:      utils.PgtypeInt4ToInt(ps.ServiceID),
		MaxRequest:     utils.PgtypeInt4ToInt(ps.MaxRequest),
		NextReset:      utils.PgtypeTimestamptzToTime(ps.NextReset),
		ResetFrequency: utils.PgtypeTextToString(ps.ResetFrequency),
		CreatedAt:      utils.PgtypeTimestamptzToTime(ps.CreatedAt),
	}
}

func ProjectServiceFromEntity(
	projectService *entities.ProjectService,
) *ProjectService {
	return &ProjectService{
		ProjectID:      utils.IntToPgtypeInt4(projectService.ProjectID),
		ServiceID:      utils.IntToPgtypeInt4(projectService.ServiceID),
		MaxRequest:     utils.IntToPgtypeInt4(projectService.MaxRequest),
		NextReset:      utils.TimeToPgtypeTimestamptz(projectService.NextReset),
		ResetFrequency: utils.StringToPgtypeText(projectService.ResetFrequency),
		CreatedAt:      utils.TimeToPgtypeTimestamptz(projectService.CreatedAt),
	}
}
