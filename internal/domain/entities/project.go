package entities

import "time"

type Project struct {
	ID int

	Name     string
	Status   string
	ClientID int

	CreatedAt time.Time
}
