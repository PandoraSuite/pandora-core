package inbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
)

type ClientHTTPPort interface {
	Create(ctx context.Context, req *dto.ClientCreate) (*dto.ClientResponse, error)
	GetClients(ctx context.Context, req *dto.ClientFilter) ([]*dto.ClientResponse, error)
}
