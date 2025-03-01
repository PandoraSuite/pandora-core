package outbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
)

type RequestLogRepositoryPort interface {
	Save(ctx context.Context, requestLog *entities.RequestLog) (*entities.RequestLog, error)
}
