package updateexecutionstatus

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type RequestRepository interface {
	UpdateExecutionStatus(ctx context.Context, id string, update *dto.RequestExecutionStatusUpdate) errors.Error
}
