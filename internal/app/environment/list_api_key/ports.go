package listapikey

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type EnvironmentRepository interface {
	Exists(ctx context.Context, id int) (bool, errors.Error)
}

type APIKeyRepository interface {
	ListByEnvironment(ctx context.Context, environmentID int) ([]*entities.APIKey, errors.Error)
}
