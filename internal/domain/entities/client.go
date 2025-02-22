package entities

import "time"

type Client struct {
	ID int

	Type  string
	Name  string
	Email string

	CreatedAt time.Time
}
