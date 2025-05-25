package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

// ... Requests ...

type ServiceFilter struct {
	Status enums.ServiceStatus `name:"status" validate:"omitempty,enums=enabled disabled deprecated"`
}

type ServiceCreate struct {
	Name    string `name:"name" validate:"required"`
	Version string `name:"version" validate:"required"`
}

// ... Responses ...

type ServiceResponse struct {
	ID        int                 `name:"id"`
	Name      string              `name:"name"`
	Status    enums.ServiceStatus `name:"status"`
	Version   string              `name:"version"`
	CreatedAt time.Time           `name:"created_at"`
}
