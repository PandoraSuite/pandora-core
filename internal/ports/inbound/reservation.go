package inbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ReservationGRPCPort interface {
	Commit(ctx context.Context, id string) *errors.Error
}
