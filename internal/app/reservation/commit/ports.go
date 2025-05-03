package commit

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ReservationRepository interface {
	Delete(ctx context.Context, id string) errors.Error
}
