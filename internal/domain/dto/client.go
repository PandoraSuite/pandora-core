package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type ClientFilter struct {
	Type enums.ClientType `json:"type,omitempty" enums:"developer,organization" swaggertype:"string"`
}

type ClientCreate struct {
	Type  enums.ClientType `json:"type" binding:"required" enums:"developer,organization" swaggertype:"string"`
	Name  string           `json:"name" binding:"required"`
	Email string           `json:"email" binding:"required,email"`
}

type ClientResponse struct {
	ID        int              `json:"id"`
	Type      enums.ClientType `json:"type" enums:"developer,organization" swaggertype:"string"`
	Name      string           `json:"name"`
	Email     string           `json:"email"`
	CreatedAt time.Time        `json:"created_at"`
}

type ClientUpdate struct {
	Type  enums.ClientType `json:"type,omitempty" enums:"developer,organization" swaggertype:"string"`
	Name  string           `json:"name,omitempty"`
	Email string           `json:"email,omitempty" binding:"omitempty,email"`
}
