package models

import "github.com/jackc/pgx/v5/pgtype"

type EnvironmentService struct {
	EnvironmentID    pgtype.Int4
	ServiceID        pgtype.Int4
	MaxRequest       pgtype.Int4
	AvailableRequest pgtype.Int4
	CreatedAt        pgtype.Timestamptz
}
