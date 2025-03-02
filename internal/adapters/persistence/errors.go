package persistence

import (
	"errors"

	domainErr "github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func ConvertPgxError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domainErr.ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "42P01":
			return domainErr.ErrUndefinedEntity
		case "23505":
			return domainErr.ErrUniqueViolation
		case "23502":
			return domainErr.ErrNotNullViolation
		case "23503":
			return domainErr.ErrForeignKeyViolation
		case "23514":
			return domainErr.ErrRestrictionViolation
		default:
			return domainErr.ErrPersistence
		}
	}
	return err
}
