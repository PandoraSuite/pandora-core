package entities

import (
	"time"

	"github.com/MAD-py/pandora-core/internal/domain/enums"
)

type Client struct {
	ID int

	Type  enums.ClientType
	Name  string
	Email string

	CreatedAt time.Time
	UpdatedAt time.Time
}
