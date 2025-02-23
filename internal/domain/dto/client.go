package dto

import "time"

type ClientType string

const (
	ClientDeveloper    ClientType = "developer"
	ClientOrganization ClientType = "organization"
)

type ClientCreate struct {
	Type  ClientType `json:"type"`
	Name  string     `json:"name"`
	Email string     `json:"email"`
}

type ClientResponse struct {
	ID        int        `json:"id"`
	Type      ClientType `json:"type"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
}
