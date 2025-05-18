package updatestatus

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ServiceRepository interface {
	UpdateStatus(ctx context.Context, id int, status enums.ServiceStatus) (*entities.Service, errors.Error)
}
