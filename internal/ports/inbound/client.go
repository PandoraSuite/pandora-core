package inbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ClientHTTPPort interface {
	Create(ctx context.Context, req *dto.ClientCreate) (*dto.ClientResponse, *errors.Error)
	GetAll(ctx context.Context, req *dto.ClientFilter) ([]*dto.ClientResponse, *errors.Error)
	Update(ctx context.Context, id int, req *dto.ClientUpdate) *errors.Error
	GetByID(ctx context.Context, id int) (*dto.ClientResponse, *errors.Error)
	GetProjects(ctx context.Context, id int) ([]*dto.ProjectResponse, *errors.Error)
}
