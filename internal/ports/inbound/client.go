package inbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/errors"
)

type ClientHTTPPort interface {
	Create(ctx context.Context, req *dto.ClientCreate) (*dto.ClientResponse, *errors.Error)
	GetClients(ctx context.Context, req *dto.ClientFilter) ([]*dto.ClientResponse, *errors.Error)
}
