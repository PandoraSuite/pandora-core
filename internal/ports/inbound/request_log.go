package inbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type RequestLogGRPCPort interface {
	UpdateExecutionStatus(ctx context.Context, id string, executionStatus enums.RequestLogExecutionStatus) *errors.Error
}
