package usecases

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
)

type APIKeyUseCase interface {
	Create(ctx context.Context, newAPIKey *entities.APIKey) error
}
