package models

import "github.com/jackc/pgx/v5/pgtype"

type EnvironmentService struct {
	ServiceID     pgtype.Int4
	EnvironmentID pgtype.Int4

	MaxRequest       pgtype.Int4
	AvailableRequest pgtype.Int4

	CreatedAt pgtype.Timestamptz
}
