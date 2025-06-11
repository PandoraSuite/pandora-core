package disabled

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type APIKeyRepository interface {
	GetByID(ctx context.Context, id int) (*entities.APIKey, errors.Error)
	UpdateStatus(ctx context.Context, id int, status enums.APIKeyStatus) errors.Error
}
