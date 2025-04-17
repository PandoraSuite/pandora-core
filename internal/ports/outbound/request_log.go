package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type RequestLogPort interface {
	Save(ctx context.Context, requestLog *entities.RequestLog) *errors.Error
	DeleteByService(ctx context.Context, serviceID int) *errors.Error
	SaveAsInitialPoint(ctx context.Context, requestLog *entities.RequestLog) *errors.Error
	UpdateExecutionStatus(ctx context.Context, id string, executionStatus enums.RequestLogExecutionStatus) *errors.Error
}
