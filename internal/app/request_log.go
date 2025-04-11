package app

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/MAD-py/pandora-core/internal/ports/outbound"
)

type RequestLogUseCase struct {
	requestLogRepo outbound.RequestLogPort
}

func (u *RequestLogUseCase) UpdateExecutionStatus(
	ctx context.Context, id string, executionStatus enums.RequestLogExecutionStatus,
) *errors.Error {
	switch executionStatus {
	case enums.RequestLogExecutionStatusNull:
		return errors.ErrCannotUpdateToNullExecutionStatus
	case enums.RequestLogPending:
		return errors.ErrCannotUpdateToPendingExecutionStatus
	case enums.RequestLogUnauthorized:
		return errors.ErrCannotUpdateToNullExecutionStatus
	default:
		return u.requestLogRepo.UpdateExecutionStatus(ctx, id, executionStatus)
	}
}

func NewRequestLogUseCase(
	requestLogRepo outbound.RequestLogPort,
) *RequestLogUseCase {
	return &RequestLogUseCase{
		requestLogRepo: requestLogRepo,
	}
}
