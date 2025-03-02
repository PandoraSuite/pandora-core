package models

import (
	"fmt"
	"slices"
	"strings"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence/models/utils"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/jackc/pgx/v5/pgtype"
)

var apiKeyStatus = []string{"active", "deactivated"}

type APIKey struct {
	ID pgtype.Int4

	Key           pgtype.Text
	Status        pgtype.Text
	LastUsed      pgtype.Timestamptz
	ExpiresAt     pgtype.Timestamptz
	EnvironmentID pgtype.Int4

	CreatedAt pgtype.Timestamptz
}

func (k *APIKey) ValidateModel() error {
	return k.validateStatus()
}

func (k *APIKey) validateStatus() error {
	if status, _ := k.Status.Value(); status != nil {
		if slices.Contains(apiKeyStatus, status.(string)) {
			return nil
		}
	}

	return fmt.Errorf(
		"invalid status: must be %s", strings.Join(apiKeyStatus, ", "),
	)
}

func (k *APIKey) ToEntity() (*entities.APIKey, error) {
	status, err := enums.ParseAPIKeyStatus(utils.PgtypeTextToString(k.Status))
	if err != nil {
		return nil, err
	}

	return &entities.APIKey{
		ID:            utils.PgtypeInt4ToInt(k.ID),
		Key:           utils.PgtypeTextToString(k.Key),
		Status:        status,
		LastUsed:      utils.PgtypeTimestamptzToTime(k.LastUsed),
		ExpiresAt:     utils.PgtypeTimestamptzToTime(k.ExpiresAt),
		EnvironmentID: utils.PgtypeInt4ToInt(k.EnvironmentID),
		CreatedAt:     utils.PgtypeTimestamptzToTime(k.CreatedAt),
	}, nil
}

func APIKeysToEntity(array []*APIKey) ([]*entities.APIKey, error) {
	result := make([]*entities.APIKey, len(array))
	for i, v := range array {
		vv, err := v.ToEntity()
		if err != nil {
			return nil, err
		}
		result[i] = vv
	}
	return result, nil
}

func APIKeyFromEntity(apiKey *entities.APIKey) *APIKey {
	return &APIKey{
		ID:            utils.IntToPgtypeInt4(apiKey.ID),
		Key:           utils.StringToPgtypeText(apiKey.Key),
		Status:        utils.StringToPgtypeText(apiKey.Status.String()),
		LastUsed:      utils.TimeToPgtypeTimestamptz(apiKey.LastUsed),
		ExpiresAt:     utils.TimeToPgtypeTimestamptz(apiKey.ExpiresAt),
		EnvironmentID: utils.IntToPgtypeInt4(apiKey.EnvironmentID),
		CreatedAt:     utils.TimeToPgtypeTimestamptz(apiKey.CreatedAt),
	}
}
