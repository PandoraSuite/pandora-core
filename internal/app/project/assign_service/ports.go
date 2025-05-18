package assignservice

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ProjectRepository interface {
	AddService(ctx context.Context, id int, service *entities.ProjectService) errors.Error
}
