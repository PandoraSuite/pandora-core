package dto

import "time"

type ClientCreate struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ClientResponse struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
