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
	Type  string `json:"type" binding:"required" enums:"developer,organization"`
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}

func (c *ClientCreate) ToDomain() *dto.ClientCreate {
	return &dto.ClientCreate{
		Type:  enums.ClientType(c.Type),
		Name:  c.Name,
		Email: c.Email,
	}
}

type ClientUpdate struct {
	Type  string `json:"type" enums:"developer,organization"`
	Name  string `json:"name"`
	Email string `json:"email"`
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
	ID        int       `json:"id"`
	Type      string    `json:"type" enums:"developer,organization"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at" time_format:"2006-01-02T15:04:05Z07:00" time_utc:"1"`
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
