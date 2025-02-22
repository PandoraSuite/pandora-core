package models

import "github.com/jackc/pgx/v5/pgtype"

type Service struct {
	ID        pgtype.Int4
	Name      pgtype.Text
	Version   pgtype.Text
	Status    pgtype.Text
	CreatedAt pgtype.Timestamptz
}
