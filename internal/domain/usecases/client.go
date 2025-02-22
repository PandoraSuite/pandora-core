package usecases

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
)

type ClientUseCase interface {
	Create(ctx context.Context, newClient *entities.Client) error
}
