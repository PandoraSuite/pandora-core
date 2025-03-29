package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type ServiceFilter struct {
	Status enums.ServiceStatus `form:"status,omitempty" enums:"active,deactivated" swaggertype:"string"`
}

type ServiceCreate struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ServiceResponse struct {
	ID        int                 `json:"id"`
	Name      string              `json:"name"`
	Status    enums.ServiceStatus `json:"status" enums:"active,deactivated" swaggertype:"string"`
	Version   string              `json:"version"`
	CreatedAt time.Time           `json:"created_at"`
}
