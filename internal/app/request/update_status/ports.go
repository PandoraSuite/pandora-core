package updatestatus

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type RequestRepository interface {
	UpdateExecutionStatus(ctx context.Context, id string, executionStatus enums.RequestExecutionStatus) errors.Error
}
