package inbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
)

type ClientPort interface {
	Create(ctx context.Context, req *dto.ClientCreate) (*dto.ClientResponse, error)
}
