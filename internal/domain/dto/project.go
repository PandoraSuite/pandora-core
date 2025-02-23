package dto

import "time"

type ProjectCreate struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	ClientID int    `json:"client_id"`
}

type ProjectResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	ClientID  int       `json:"client_id"`
	CreatedAt time.Time `json:"created_at"`
}
