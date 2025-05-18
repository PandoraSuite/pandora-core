package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

// ... Requests ...

type ClientFilter struct {
	Type enums.ClientType `name:"type" validate:"omitempty,enums=developer organization"`
}

type ClientCreate struct {
	Type  enums.ClientType `name:"type" validate:"required,enums=developer organization"`
	Name  string           `name:"name" validate:"required"`
	Email string           `name:"email" validate:"required,email"`
}

type ClientUpdate struct {
	Type  enums.ClientType `name:"type" validate:"omitempty,enums=developer organization"`
	Name  string           `name:"name" validate:"omitempty"`
	Email string           `name:"email" validate:"omitempty,email"`
}

// ... Responses ...

type ClientResponse struct {
	ID        int              `name:"id"`
	Type      enums.ClientType `name:"type"`
	Name      string           `name:"name"`
	Email     string           `name:"email"`
	CreatedAt time.Time        `name:"created_at"`
}
