package models

import (
	"fmt"
	"slices"
	"strings"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence/models/utils"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/jackc/pgx/v5/pgtype"
)

var requestLogExecutionStatus = []string{
	"success", "failed", "unauthorized", "server error",
}

type RequestLog struct {
	ID pgtype.Int4

	APIKey          pgtype.Text
	ServiceID       pgtype.Int4
	RequestTime     pgtype.Timestamptz
	EnvironmentID   pgtype.Int4
	ExecutionStatus pgtype.Text

	CreatedAt pgtype.Timestamptz
}

func (rl *RequestLog) ValidateModel() error {
	return rl.validateExecutionStatus()
}

func (rl *RequestLog) validateExecutionStatus() error {
	if executionStatus, _ := rl.ExecutionStatus.Value(); executionStatus != nil {
		if slices.Contains(requestLogExecutionStatus, executionStatus.(string)) {
			return nil
		}
	}

	return fmt.Errorf(
		"invalid execution status: must be %s",
		strings.Join(requestLogExecutionStatus, ", "),
	)
}

func (rl *RequestLog) ToEntity() *entities.RequestLog {
	return &entities.RequestLog{
		ID:              utils.PgtypeInt4ToInt(rl.ID),
		APIKey:          utils.PgtypeTextToString(rl.APIKey),
		ServiceID:       utils.PgtypeInt4ToInt(rl.ServiceID),
		RequestTime:     utils.PgtypeTimestamptzToTime(rl.RequestTime),
		EnvironmentID:   utils.PgtypeInt4ToInt(rl.EnvironmentID),
		ExecutionStatus: utils.PgtypeTextToString(rl.ExecutionStatus),
		CreatedAt:       utils.PgtypeTimestamptzToTime(rl.CreatedAt),
	}
}

func RequestLogFromEntity(requestLog *entities.RequestLog) *RequestLog {
	return &RequestLog{
		ID:              utils.IntToPgtypeInt4(requestLog.ID),
		APIKey:          utils.StringToPgtypeText(requestLog.APIKey),
		ServiceID:       utils.IntToPgtypeInt4(requestLog.ServiceID),
		RequestTime:     utils.TimeToPgtypeTimestamptz(requestLog.RequestTime),
		EnvironmentID:   utils.IntToPgtypeInt4(requestLog.EnvironmentID),
		ExecutionStatus: utils.StringToPgtypeText(requestLog.ExecutionStatus),
		CreatedAt:       utils.TimeToPgtypeTimestamptz(requestLog.CreatedAt),
	}
}
