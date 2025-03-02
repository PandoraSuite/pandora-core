package persistence

import (
	"errors"

	persistenceErr "github.com/MAD-py/pandora-core/internal/domain/errors"
	"github.com/jackc/pgx/v5/pgconn"
)

func ConvertPgxError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "42P01":
			return persistenceErr.ErrUndefinedEntity
		case "23505":
			return persistenceErr.ErrUniqueViolation
		case "23503":
			return persistenceErr.ErrForeignKeyViolation
		case "23514":
			return persistenceErr.ErrRestrictionViolation
		default:
			return persistenceErr.ErrPersistence
		}
	}
	return err
}
