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

var serviceStatus = []string{"active", "deactivated", "deprecated"}

type Service struct {
	ID pgtype.Int4

	Name    pgtype.Text
	Status  pgtype.Text
	Version pgtype.Text

	CreatedAt pgtype.Timestamptz
}

func (s *Service) EntityID() int {
	return utils.PgtypeInt4ToInt(s.ID)
}

func (s *Service) EntityCreatedAt() time.Time {
	return utils.PgtypeTimestamptzToTime(s.CreatedAt)
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

func (s *Service) ToEntity() (*entities.Service, error) {
	status, err := enums.ParseServiceStatus(utils.PgtypeTextToString(s.Status))
	if err != nil {
		return nil, err
	}

	return &entities.Service{
		ID:        utils.PgtypeInt4ToInt(s.ID),
		Name:      utils.PgtypeTextToString(s.Name),
		Status:    status,
		Version:   utils.PgtypeTextToString(s.Version),
		CreatedAt: utils.PgtypeTimestamptzToTime(s.CreatedAt),
	}, nil
}

func ServicesToEntity(array []*Service) ([]*entities.Service, error) {
	result := make([]*entities.Service, len(array))
	for i, v := range array {
		vv, err := v.ToEntity()
		if err != nil {
			return nil, err
		}
		result[i] = vv
	}
	return result, nil
}

func ServiceFromEntity(service *entities.Service) Service {
	return Service{
		ID:        utils.IntToPgtypeInt4(service.ID),
		Name:      utils.StringToPgtypeText(service.Name),
		Status:    utils.StringToPgtypeText(service.Status.String()),
		Version:   utils.StringToPgtypeText(service.Version),
		CreatedAt: utils.TimeToPgtypeTimestamptz(service.CreatedAt),
	}
}
