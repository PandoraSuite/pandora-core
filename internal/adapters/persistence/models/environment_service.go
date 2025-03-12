package models

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence/models/utils"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/jackc/pgx/v5/pgtype"
)

type EnvironmentService struct {
	ServiceID     pgtype.Int4
	EnvironmentID pgtype.Int4

	MaxRequest       pgtype.Int4
	AvailableRequest pgtype.Int4

	CreatedAt pgtype.Timestamptz
}

func (es *EnvironmentService) EntityCreatedAt() time.Time {
	return utils.PgtypeTimestamptzToTime(es.CreatedAt)
}

func (es *EnvironmentService) ToEntity() *entities.EnvironmentService {
	return &entities.EnvironmentService{
		ServiceID:        utils.PgtypeInt4ToInt(es.ServiceID),
		EnvironmentID:    utils.PgtypeInt4ToInt(es.EnvironmentID),
		MaxRequest:       utils.PgtypeInt4ToInt(es.MaxRequest),
		AvailableRequest: utils.PgtypeInt4ToInt(es.AvailableRequest),
		CreatedAt:        utils.PgtypeTimestamptzToTime(es.CreatedAt),
	}
}

func EnvironmentServicesToEntity(
	array []*EnvironmentService,
) ([]*entities.EnvironmentService, error) {
	result := make([]*entities.EnvironmentService, len(array))
	for i, v := range array {
		result[i] = v.ToEntity()
	}
	return result, nil
}

func EnvironmentServiceFromEntity(
	environmentService *entities.EnvironmentService,
) EnvironmentService {
	return EnvironmentService{
		ServiceID:        utils.IntToPgtypeInt4(environmentService.ServiceID),
		EnvironmentID:    utils.IntToPgtypeInt4(environmentService.EnvironmentID),
		MaxRequest:       utils.IntToPgtypeInt4(environmentService.MaxRequest),
		AvailableRequest: utils.IntToPgtypeInt4(environmentService.AvailableRequest),
		CreatedAt:        utils.TimeToPgtypeTimestamptz(environmentService.CreatedAt),
	}
}
