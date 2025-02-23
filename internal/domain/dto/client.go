package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type ClientCreate struct {
	Type  enums.ClientType `json:"type"`
	Name  string           `json:"name"`
	Email string           `json:"email"`
}

type ClientResponse struct {
	ID        int              `json:"id"`
	Type      enums.ClientType `json:"type"`
	Name      string           `json:"name"`
	Email     string           `json:"email"`
	CreatedAt time.Time        `json:"created_at"`
}
