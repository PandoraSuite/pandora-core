package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	domainErr "github.com/MAD-py/pandora-core/internal/domain/errors"
)

type Driver struct {
	pool *pgxpool.Pool
}

func (d *Driver) Close() {
	if d.pool != nil {
		d.pool.Close()
	}
}

func (d *Driver) entityMapper(table string) string {
	switch table {
	case "service":
		return "Service"
	case "client":
		return "Client"
	case "project":
		return "Project"
	case "environment":
		return "Environment"
	case "api_key":
		return "APIKey"
	case "project_service":
		return "ProjectService"
	case "environment_service":
		return "EnvironmentService"
	case "request_log":
		return "RequestLog"
	case "reservation":
		return "Reservation"
	default:
		return table
	}
}

func (d *Driver) errorMapper(err error, tableName string) domainErr.PersistenceError {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return domainErr.NewPersistenceNotFoundError(d.entityMapper(tableName), err)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // UNIQUE_VIOLATION
			return domainErr.NewPersistenceError(
				domainErr.PersistenceErrorCodeUniqueViolation,
				d.entityMapper(pgErr.TableName),
				pgErr.ColumnName,
				pgErr.Message,
				pgErr,
			)
		case "23503": // FOREIGN_KEY_VIOLATION
			return domainErr.NewPersistenceError(
				domainErr.PersistenceErrorCodeInvalidReference,
				d.entityMapper(pgErr.TableName),
				pgErr.ColumnName,
				pgErr.Message,
				pgErr,
			)
		case "23502", // NOT_NULL_VIOLATION
			"23514": // CHECK_VIOLATION
			return domainErr.NewPersistenceError(
				domainErr.PersistenceErrorCodeInvalidValue,
				d.entityMapper(pgErr.TableName),
				pgErr.ColumnName,
				pgErr.Message,
				pgErr,
			)
		case "42P01": // UNDEFINED_TABLE
			return domainErr.NewPersistenceError(
				domainErr.PersistenceErrorCodeUndefinedEntity,
				d.entityMapper(pgErr.TableName),
				pgErr.ColumnName,
				pgErr.Message,
				pgErr,
			)
		case "08001", // SQLCLIENT_UNABLE_TO_ESTABLISH_SQLCONNECTION
			"08006", // CONNECTION_FAILURE
			"08003", // CONNECTION_DOES_NOT_EXIST
			"57P01": // ADMIN_SHUTDOWN
			return domainErr.NewPersistenceConnectionError(pgErr.Message, pgErr)
		default:
			return domainErr.NewPersistenceUnknownError(pgErr.Message, pgErr)
		}
	}

	return domainErr.NewPersistenceUnknownError(err.Error(), err)
}

func (d *Driver) entityNotFoundError(tableName string) domainErr.PersistenceError {
	return domainErr.NewPersistenceNotFoundError(d.entityMapper(tableName), nil)
}

func NewDriver(dns string) *Driver {
	config, err := pgxpool.ParseConfig(dns)
	if err != nil {
		panic(err)
	}

	config.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		panic(err)
	}

	return &Driver{pool: pool}
}
