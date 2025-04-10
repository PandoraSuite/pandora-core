package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type RequestLogPort interface {
	Save(ctx context.Context, requestLog *entities.RequestLog) *errors.Error
	SaveAsInitialPoint(ctx context.Context, requestLog *entities.RequestLog) *errors.Error
}
