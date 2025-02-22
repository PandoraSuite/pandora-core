package entities

import "time"

type Environment struct {
	ID int

	Name      string
	Status    string
	ProjectID int

	CreatedAt time.Time
}
