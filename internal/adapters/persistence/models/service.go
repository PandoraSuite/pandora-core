package models

import (
	"fmt"
	"slices"
	"strings"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence/models/utils"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/jackc/pgx/v5/pgtype"
)

var serviceStatus = []string{"active", "deactivated", "deprecated"}

type Service struct {
	ID pgtype.Int4

	Name    pgtype.Text
	Status  pgtype.Text
	Version pgtype.Text

	CreatedAt pgtype.Timestamptz
}

func (s *Service) ValidateModel() error {
	return s.validateStatus()
}

func (s *Service) validateStatus() error {
	if status, _ := s.Status.Value(); status != nil {
		if slices.Contains(serviceStatus, status.(string)) {
			return nil
		}
	}

	return fmt.Errorf(
		"invalid status: must be %s", strings.Join(serviceStatus, ", "),
	)
}

func (s *Service) ToEntity() *entities.Service {
	return &entities.Service{
		ID:        utils.PgtypeInt4ToInt(s.ID),
		Name:      utils.PgtypeTextToString(s.Name),
		Status:    utils.PgtypeTextToString(s.Status),
		Version:   utils.PgtypeTextToString(s.Version),
		CreatedAt: utils.PgtypeTimestamptzToTime(s.CreatedAt),
	}
}

func ServiceFromEntity(service *entities.Service) *Service {
	return &Service{
		ID:        utils.IntToPgtypeInt4(service.ID),
		Name:      utils.StringToPgtypeText(service.Name),
		Status:    utils.StringToPgtypeText(service.Status),
		Version:   utils.StringToPgtypeText(service.Version),
		CreatedAt: utils.TimeToPgtypeTimestamptz(service.CreatedAt),
	}
}
