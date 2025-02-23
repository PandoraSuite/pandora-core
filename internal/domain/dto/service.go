package dto

import "time"

type ServiceCreate struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Version string `json:"version"`
}

type ServiceResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Version   string    `json:"version"`
	CreatedAt time.Time `json:"created_at"`
}
