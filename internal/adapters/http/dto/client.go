package dto

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/dto"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

// ... Requests ...

type ClientFilter struct {
	Type string `form:"type" enums:"developer,organization"`
}

func (c *ClientFilter) ToDomain() *dto.ClientFilter {
	return &dto.ClientFilter{
		Type: enums.ClientType(c.Type),
	}
}

type ClientCreate struct {
	Type string `json:"type" validate:"required" enums:"developer,organization"`

	Name string `json:"name" validate:"required"`

	Email string `json:"email" validate:"required" format:"email"`
}

func (c *ClientCreate) ToDomain() *dto.ClientCreate {
	return &dto.ClientCreate{
		Type:  enums.ClientType(c.Type),
		Name:  c.Name,
		Email: c.Email,
	}
}

type ClientUpdate struct {
	Type string `json:"type" enums:"developer,organization"`

	Name string `json:"name"`

	Email string `json:"email" format:"email"`
}

func (c *ClientUpdate) ToDomain() *dto.ClientUpdate {
	return &dto.ClientUpdate{
		Type:  enums.ClientType(c.Type),
		Name:  c.Name,
		Email: c.Email,
	}
}

// ... Responses ...

type ClientResponse struct {
	ID int `json:"id" minimum:"1"`

	Type string `json:"type" enums:"developer,organization"`

	Name string `json:"name"`

	Email string `json:"email" format:"email"`

	CreatedAt time.Time `json:"created_at" format:"date-time" extensions:"x-timezone=utc"`
}

func ClientResponseFromDomain(client *dto.ClientResponse) *ClientResponse {
	return &ClientResponse{
		ID:        client.ID,
		Type:      string(client.Type),
		Name:      client.Name,
		Email:     client.Email,
		CreatedAt: client.CreatedAt,
	}
}
