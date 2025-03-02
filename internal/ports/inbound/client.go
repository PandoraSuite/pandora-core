package inbound

import (
	"context"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type ClientPort interface {
	Create(ctx context.Context, req *dto.ClientCreate) (*dto.ClientResponse, error)
	GetClients(ctx context.Context, clientType enums.ClientType) ([]*dto.ClientResponse, error)
}
