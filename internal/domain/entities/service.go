package entities

import "time"

type Service struct {
	ID int

	Name    string
	Status  string
	Version string

	CreatedAt time.Time
}
