package persistence

import (
	"context"
	"errors"
	"time"

	domainErr "github.com/MAD-py/pandora-core/internal/domain/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Persistence struct {
	pool *pgxpool.Pool
}

func (db *Persistence) Close() { db.pool.Close() }

func (db *Persistence) Pool() *pgxpool.Pool { return db.pool }

func (db *Persistence) HandlerErr() func(error) error {
	return func(err error) error {
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
}

func NewPersistence(dns string) (*Persistence, error) {
	config, err := pgxpool.ParseConfig(dns)
	if err != nil {
		return nil, err
	}

	config.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return &Persistence{pool: pool}, nil
}
